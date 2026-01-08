package components

import (
	"fmt"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/internal"
	. "github.com/AnatoleLucet/loom/components"
)

func FullscreenRenderer(fn func() loom.Node) loom.Node {
	return &fullscreenRendererNode{
		fn: fn,
	}
}

type fullscreenRendererNode struct {
	fn func() loom.Node
}

func (n *fullscreenRendererNode) ID() string {
	return "term.FullscreenRenderer"
}

func (n *fullscreenRendererNode) Mount(slot *loom.Slot) error {
	ctx, err := NewAppContext(internal.RendererTypeFullscreen)
	if err != nil {
		return fmt.Errorf("FullscreenRenderer: %w", err)
	}

	slot.SetSelf(ctx.renderer.Elem())

	return ctx.renderer.Batch(func() error {
		return slot.RenderChildren(Provider(appctx, ctx, n.fn))
	})
}

func (n *fullscreenRendererNode) Update(slot *loom.Slot) error {
	ctx, err := getAppContext()
	if err != nil {
		return fmt.Errorf("FullscreenRenderer: %w. That's not supposed to happend. Please repport an issue at https://github.com/AnatoleLucet/loom/issues", err)
	}

	return ctx.renderer.Batch(func() error {
		return slot.RenderChildren(Provider(appctx, ctx, n.fn))
	})
}

func (n *fullscreenRendererNode) Unmount(slot *loom.Slot) error {
	ctx, err := getAppContext()
	if err != nil {
		return fmt.Errorf("FullscreenRenderer: %w. That's not supposed to happend. Please repport an issue at https://github.com/AnatoleLucet/loom/issues", err)
	}

	ctx.renderer.Close()
	return nil
}

func InlineRenderer(fn func() loom.Node) loom.Node {
	return &inlineRendererNode{
		fn: fn,
	}
}

type inlineRendererNode struct {
	fn func() loom.Node
}

func (n *inlineRendererNode) ID() string {
	return "term.InlineRenderer"
}

func (n *inlineRendererNode) Mount(slot *loom.Slot) error {
	ctx, err := NewAppContext(internal.RendererTypeInline)
	if err != nil {
		return fmt.Errorf("InlineRenderer: %w", err)
	}

	slot.SetSelf(ctx.renderer.Elem())

	return ctx.renderer.Batch(func() error {
		return slot.RenderChildren(Provider(appctx, ctx, n.fn))
	})
}

func (n *inlineRendererNode) Update(slot *loom.Slot) error {
	// todo: that's not going to work (same for unmount and fullscreen)
	ctx, err := getAppContext()
	if err != nil {
		return fmt.Errorf("InlineRenderer: %w. That's not supposed to happend. Please repport an issue at https://github.com/AnatoleLucet/loom/issues", err)
	}

	return ctx.renderer.Batch(func() error {
		return slot.RenderChildren(Provider(appctx, ctx, n.fn))
	})
}

func (n *inlineRendererNode) Unmount(slot *loom.Slot) error {
	ctx, err := getAppContext()
	if err != nil {
		return fmt.Errorf("InlineRenderer: %w. That's not supposed to happend. Please repport an issue at https://github.com/AnatoleLucet/loom/issues", err)
	}

	ctx.renderer.Close()
	return nil
}
