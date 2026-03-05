package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/components/context"
	"github.com/AnatoleLucet/loom-term/core"
)

func TextArea(children ...loom.Node) loom.Node {
	return &textAreaNode{
		children: children,
	}
}

type textAreaNode struct {
	children []loom.Node
}

func (n *textAreaNode) ID() string {
	return "term.TextArea"
}

func (n *textAreaNode) Mount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("TextArea: %w", err)
	}

	parent := slot.Parent().(core.Element)
	self, err := core.NewTextAreaElement()
	if err != nil {
		return fmt.Errorf("TextArea: %w", err)
	}
	slot.SetSelf(self)

	return ctx.BatchRender(func() error {
		parent.AppendChild(self)
		return n.Update(slot)
	})
}

func (n *textAreaNode) Update(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("TextArea: %w", err)
	}

	return ctx.BatchRender(func() error {
		return slot.RenderChildren(n.children...)
	})
}

func (n *textAreaNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("TextArea: %w", err)
	}

	self := slot.Self().(core.Element)

	return ctx.BatchRender(func() error {
		self.Destroy()
		return nil
	})
}
