package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/app"
	. "github.com/AnatoleLucet/loom/components"
)

func Text(content any, styles ...*applyNode) loom.Node {
	children := make([]loom.Node, len(styles))
	for i, style := range styles {
		children[i] = style
	}

	return &textNode{
		name:     "Text",
		content:  fmt.Sprint(content),
		children: children,
	}
}
func BindText[T any](fn func() T, styles ...*applyNode) loom.Node {
	return Bind(func() loom.Node {
		return Text(fn(), styles...)
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

func Br() loom.Node {
	return &textNode{
		name: "Br",
		modifier: func(parent core.Element, self core.TextElement) {
			_, ok := parent.(core.TextElement)
			if ok {
				self.SetText("\n")
			}
		},
	}
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
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("P: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	parent := slot.Parent().(core.Element)
	self, err := core.NewTextElement()
	if err != nil {
		return fmt.Errorf("P: %w", err)
	}
	slot.SetSelf(self)

	if n.modifier != nil {
		n.modifier(parent, self)
	}

	parent.AppendChild(self)
	ctx.RequestRender()

	return n.Update(slot)
}

func (n *textNode) Update(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("P: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	self := slot.Self().(core.TextElement)
	if n.content != "" {
		self.SetText(n.content)
	}

	return slot.RenderChildren(n.children...)
}

func (n *textNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("P: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	self := slot.Self().(core.Element)

	self.Destroy()
	ctx.RequestRender()
	return nil
}
