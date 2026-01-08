package components

import (
	"fmt"
	"math"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/internal"
	"github.com/AnatoleLucet/loom-term/opentui"
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
	parent := slot.Parent().(*internal.Element)

	self, err := internal.NewElement(parent)
	if err != nil {
		return err
	}
	slot.SetSelf(self)
	parent.AppendChild(self)

	self.SetPaint(n.Paint)

	return n.Update(slot)
}

func (n *boxNode) Update(slot *loom.Slot) error {
	ctx, err := getAppContext()
	if err != nil {
		return fmt.Errorf("Box: %w", err)
	}

	return ctx.renderer.Batch(func() error {
		return slot.RenderChildren(n.children...)
	})
}

func (n *boxNode) Unmount(slot *loom.Slot) error {
	parent := slot.Parent().(*internal.Element)
	self := slot.Self().(*internal.Element)

	parent.RemoveChild(self)

	return n.Update(slot)
}

func (n *boxNode) Paint(node *internal.Element, buffer *opentui.Buffer) error {
	layout := node.Layout().GetLayout()

	// Calculate absolute position by walking up the tree
	left := layout.Left()
	top := layout.Top()
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		parentLayout := parent.Layout().GetLayout()
		left += parentLayout.Left()
		top += parentLayout.Top()
	}

	// Use rounding instead of truncation for consistent sub-pixel positioning
	buffer.FillRect(
		uint32(math.Round(float64(left))),
		uint32(math.Round(float64(top))),
		uint32(math.Round(float64(layout.Width()))),
		uint32(math.Round(float64(layout.Height()))),
		node.BackgroundColor(),
	)

	return nil
}
