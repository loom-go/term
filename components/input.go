package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/components/context"
	"github.com/AnatoleLucet/loom-term/core"
)

func Input(children ...loom.Node) loom.Node {
	return &inputNode{
		children: children,
	}
}

type inputNode struct {
	children []loom.Node
}

func (n *inputNode) ID() string {
	return "term.Input"
}

func (n *inputNode) Mount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Input: %w", err)
	}

	parent := slot.Parent().(core.Element)
	self, err := core.NewInputElement()
	if err != nil {
		return fmt.Errorf("Input: %w", err)
	}
	slot.SetSelf(self)

	return ctx.BatchRender(func() error {
		parent.AppendChild(self)
		return n.Update(slot)
	})
}

func (n *inputNode) Update(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Input: %w", err)
	}

	return ctx.BatchRender(func() error {
		return slot.RenderChildren(n.children...)
	})
}

func (n *inputNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("Input: %w", err)
	}

	self := slot.Self().(core.Element)

	return ctx.BatchRender(func() error {
		self.Destroy()
		return nil
	})
}
