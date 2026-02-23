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

func (n *consoleNode) ID() string {
	return "term.Console"
}

func (n *consoleNode) Mount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	parent := slot.Parent().(core.Element)
	self, err := core.NewConsoleElement()
	if err != nil {
		return err
	}
	slot.SetSelf(self)

	parent.AppendChild(self)
	ctx.ScheduleRender()

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

	self := slot.Self().(core.Element)

	self.Destroy()
	ctx.ScheduleRender()

	return nil
}
