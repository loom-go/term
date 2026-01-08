package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/internal"
	"github.com/AnatoleLucet/loom-term/opentui"
	. "github.com/AnatoleLucet/loom/components"
)

func Text(content string) loom.Node {
	return &textNode{content}
}

func BindText(fn func() string) loom.Node {
	return Bind(func() loom.Node {
		return Text(fn())
	})
}

type textNode struct {
	content string
}

func (n *textNode) ID() string {
	return "term.Text"
}

func (n *textNode) Mount(slot *loom.Slot) error {
	parent := slot.Parent().(*internal.Element)

	self, err := internal.NewElement(parent)
	if err != nil {
		return err
	}
	slot.SetSelf(self)
	parent.AppendChild(self)

	return n.Update(slot)
}

func (n *textNode) Update(slot *loom.Slot) error {
	ctx, err := getAppContext()
	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	self := slot.Self().(*internal.Element)
	self.SetText(n.content)
	self.SetTextColor(opentui.White)
	self.SetPaint(n.Paint)

	return ctx.renderer.Batch(func() error { return nil })
}

func (n *textNode) Unmount(slot *loom.Slot) error {
	parent := slot.Parent().(*internal.Element)
	self := slot.Self().(*internal.Element)

	parent.RemoveChild(self)

	return n.Update(slot)
}

func (n *textNode) Paint(node *internal.Element, buffer *opentui.Buffer) error {
	layout := node.Layout().GetLayout()

	// todo: helper method in tess Layout.AbsoluteLeft/Top()
	left := layout.Left()
	top := layout.Top()
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		parentLayout := parent.Layout().GetLayout()
		left += parentLayout.Left()
		top += parentLayout.Top()
	}

	buffer.DrawText(
		node.Text(),
		uint32(left),
		uint32(top),
		node.TextColor(),
		&opentui.Transparent,
		0,
	)

	return nil
}
