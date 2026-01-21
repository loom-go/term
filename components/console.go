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
	self, err := core.NewConsoleElement(ctx.RenderContext())
	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}
	slot.SetSelf(self)

	err = ctx.DoSafely(func() error {
		err = parent.AppendChild(self)
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}

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

	parent := slot.Parent().(core.Element)
	self := slot.Self().(core.Element)

	err = ctx.DoSafely(func() error {
		err = parent.RemoveChild(self)
		err = self.Destroy()
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Console: %w", err)
	}

	return nil
}
