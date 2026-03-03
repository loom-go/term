package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/core/elements"
	"github.com/AnatoleLucet/loom-term/internal/app"
)

func Box(children ...loom.Node) loom.Node {
	return &boxNode{
		children: children,
	}
}

type boxNode struct {
	children []loom.Node
}

func (n *boxNode) ID() string {
	return "term.Box"
}

func (n *boxNode) Mount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	parent := slot.Parent().(core.Element)
	self, err := elements.NewBoxElement()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}
	slot.SetSelf(self)

	parent.AppendChild(self)
	ctx.RequestRender()

	return n.Update(slot)
}

func (n *boxNode) Update(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return err
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	return slot.RenderChildren(n.children...)

}

func (n *boxNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	self := slot.Self().(core.Element)

	self.Destroy()
	ctx.RequestRender()

	return nil
}
