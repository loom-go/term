package components

import (
	"fmt"

	"github.com/loom-go/loom"
	"github.com/loom-go/term/components/appctx"
	"github.com/loom-go/term/core"
)

func ScrollBox(children ...loom.Node) loom.Node {
	return &scrollBoxNode{
		children: children,
	}
}

type scrollBoxNode struct {
	children []loom.Node
}

func (n *scrollBoxNode) ID() string {
	return "term.ScrollBox"
}

func (n *scrollBoxNode) Mount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("ScrollBox: %w", err)
	}

	parent := slot.Parent().(core.Element)
	self, err := core.NewScrollBoxElement()
	if err != nil {
		return fmt.Errorf("ScrollBox: %w", err)
	}
	slot.SetSelf(self)

	return ctx.BatchRender(func() error {
		parent.AppendChild(self)
		return n.Update(slot)
	})
}

func (n *scrollBoxNode) Update(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("ScrollBox: %w", err)
	}

	return ctx.BatchRender(func() error {
		return slot.RenderChildren(n.children...)
	})
}

func (n *scrollBoxNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("ScrollBox: %w", err)
	}

	self := slot.Self().(core.Element)

	return ctx.BatchRender(func() error {
		self.Destroy()
		return nil
	})
}
