package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/core/elements"
	"github.com/AnatoleLucet/loom-term/internal/app"
	. "github.com/AnatoleLucet/loom/components"
)

func Text(content any, styles ...loom.Node) loom.Node {
	return &textNode{
		content: fmt.Sprintf("%v", content),
		styles:  styles,
	}
}

func BindText[T any](fn func() T, styles ...loom.Node) loom.Node {
	return Bind(func() loom.Node {
		return Text(fn(), styles...)
	})
}

type textNode struct {
	content string
	styles  []loom.Node
}

func (n *textNode) ID() string {
	return "term.Text"
}

func (n *textNode) Mount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	parent := slot.Parent().(core.Element)
	self, err := elements.NewTextElement()
	if err != nil {
		return err
	}
	slot.SetSelf(self)

	parent.AppendChild(self)
	ctx.ScheduleRender()

	return n.Update(slot)
}

func (n *textNode) Update(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	self := slot.Self().(*elements.TextElement)

	self.SetText(n.content)
	ctx.ScheduleRender()

	return slot.RenderChildren(n.styles...)
}

func (n *textNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	self := slot.Self().(core.Element)

	self.Destroy()
	ctx.ScheduleRender()

	return nil
}
