package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	appctx "github.com/AnatoleLucet/loom-term/components/context"
	"github.com/AnatoleLucet/loom-term/core"
)

func Console(enabled ...bool) loom.Node {
	if len(enabled) > 0 && !enabled[0] {
		return nil
	}

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
	ctx, err := appctx.Get()
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
			event.PreventDefault()

			if self.visible {
				ctx.Root().RemoveChild(self.element)
				self.visible = false
			} else {
				ctx.Root().AppendChild(self.element)
				self.visible = true
			}
		}
	})

	return nil
}

func (n *consoleNode) Update(slot *loom.Slot) error {
	return nil
}

func (n *consoleNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}

	self := slot.Self().(*consoleState)

	return ctx.BatchRender(func() error {
		self.element.Destroy()
		if self.remove != nil {
			self.remove()
		}

		return nil
	})
}
