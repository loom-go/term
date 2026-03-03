package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/app"
)

func Console() loom.Node {
	return &consoleNode{}
}

type consoleNode struct{}

type consoleState struct {
	visible bool
	remove  func()
	element core.ConsoleElement
}

func (n *consoleNode) ID() string {
	return "term.Console"
}

func (n *consoleNode) Mount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}

	self := &consoleState{}
	slot.SetSelf(self)

	self.element, err = core.NewConsoleElement()
	if err != nil {
		return err
	}

	self.remove = ctx.Root().OnKeyPress(func(event *core.EventKey) {
		if event.Key.String() == "`" {
			if self.visible {
				ctx.Root().RemoveChild(self.element)
				self.visible = false
			} else {
				ctx.Root().AppendChild(self.element)
				self.visible = true
			}

			ctx.ScheduleRender()
		}
	})

	return nil
}

func (n *consoleNode) Update(slot *loom.Slot) error {
	return nil
}

func (n *consoleNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	self := slot.Self().(*consoleState)

	self.element.Destroy()
	if self.remove != nil {
		self.remove()
	}

	ctx.RequestRender()

	return nil
}
