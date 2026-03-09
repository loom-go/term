package components

import (
	"fmt"

	"github.com/loom-go/loom"
	. "github.com/loom-go/loom/components"
	"github.com/loom-go/term/components/appctx"
	"github.com/loom-go/term/core"
)

func Text(content any, children ...loom.Node) loom.Node {
	return &textNode{
		name:     "Text",
		content:  fmt.Sprint(content),
		children: children,
	}
}
func BindText[T any](fn func() T, children ...loom.Node) loom.Node {
	return Bind(func() loom.Node {
		return Text(fn(), children...)
	})
}

func P(children ...loom.Node) loom.Node {
	return &textNode{
		name:     "P",
		children: children,
	}
}
func Paragraph(children ...loom.Node) loom.Node {
	return P(children...)
}

func B(children ...loom.Node) loom.Node {
	return &textNode{
		name:     "B",
		children: children,
		modifier: func(parent core.Element, self core.TextElement) {
			self.SetFontWeight("bold")
		},
	}
}
func Bold(children ...loom.Node) loom.Node {
	return B(children...)
}

func I(children ...loom.Node) loom.Node {
	return &textNode{
		name:     "I",
		children: children,
		modifier: func(parent core.Element, self core.TextElement) {
			self.SetFontStyle("italic")
		},
	}
}
func Italic(children ...loom.Node) loom.Node {
	return I(children...)
}

func U(children ...loom.Node) loom.Node {
	return &textNode{
		name: "U",
		modifier: func(parent core.Element, self core.TextElement) {
			self.SetTextDecoration("underline")
		},
		children: children,
	}
}
func Underline(children ...loom.Node) loom.Node {
	return U(children...)
}

func S(children ...loom.Node) loom.Node {
	return &textNode{
		name: "S",
		modifier: func(parent core.Element, self core.TextElement) {
			self.SetTextDecoration("strikethrough")
		},
		children: children,
	}
}
func Strikethrough(children ...loom.Node) loom.Node {
	return S(children...)
}

type textNode struct {
	name     string
	content  string
	modifier func(parent core.Element, self core.TextElement)
	children []loom.Node
}

func (n *textNode) ID() string {
	return fmt.Sprintf("term.%s", n.name)
}

func (n *textNode) Mount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("%s: %w", n.name, err)
	}

	parent := slot.Parent().(core.Element)
	self, err := core.NewTextElement()
	if err != nil {
		return fmt.Errorf("%s: %w", n.name, err)
	}
	slot.SetSelf(self)

	return ctx.BatchRender(func() error {
		if n.modifier != nil {
			n.modifier(parent, self)
		}

		parent.AppendChild(self)

		return n.Update(slot)
	})
}

func (n *textNode) Update(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("%s: %w", n.name, err)
	}

	self := slot.Self().(core.TextElement)

	return ctx.BatchRender(func() error {
		if n.content != "" {
			self.SetText(n.content)
		}

		return slot.RenderChildren(n.children...)
	})
}

func (n *textNode) Unmount(slot *loom.Slot) error {
	ctx, err := appctx.Get()
	if err != nil {
		return fmt.Errorf("%s: %w", n.name, err)
	}

	self := slot.Self().(core.TextElement)

	return ctx.BatchRender(func() error {
		self.Destroy()
		return nil
	})
}
