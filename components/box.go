package components

import (
	"fmt"

	"github.com/loom-go/loom"
	"github.com/loom-go/term/components/appctx"
	"github.com/loom-go/term/core"
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
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	parent := slot.Parent().(core.Element)
	self, err := core.NewBoxElement()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}
	slot.SetSelf(self)

	return ctx.BatchRender(func() error {
		parent.AppendChild(self)
		return n.Update(slot)
	})
}

func (n *boxNode) Update(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	return ctx.BatchRender(func() error {
		return slot.RenderChildren(n.children...)
	})
}

func (n *boxNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	self := slot.Self().(core.Element)

	return ctx.BatchRender(func() error {
		self.Destroy()
		return nil
	})
}
