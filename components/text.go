package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/core/elements"
	"github.com/AnatoleLucet/loom-term/internal/app"
	. "github.com/AnatoleLucet/loom/components"
)

func Text(content string, styles ...loom.Node) loom.Node {
	return &textNode{content, styles}
}

func BindText(fn func() string, styles ...loom.Node) loom.Node {
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
	self, err := elements.NewTextElement(ctx.RenderContext())
	if err != nil {
		return err
	}
	slot.SetSelf(self)

	err = ctx.DoSafely(func() error {
		err = parent.AppendChild(self)
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

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

	err = ctx.DoSafely(func() error {
		err = self.SetContent(n.content)
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	return slot.RenderChildren(n.styles...)
}

func (n *textNode) Unmount(slot *loom.Slot) error {
	ctx, err := app.GetContext()
	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	ctx.PushRenderHold()
	defer ctx.PopRenderHold()

	parent := slot.Parent().(core.Element)
	self := slot.Self().(core.Element)

	ctx.DoSafely(func() error {
		err = parent.RemoveChild(self)
		err = self.Destroy()
		if err != nil {
			return err
		}

		return ctx.RequestRender()
	})

	if err != nil {
		return fmt.Errorf("Text: %w", err)
	}

	return nil
}
