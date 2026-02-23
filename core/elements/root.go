package elements

import (
	"context"
	"fmt"

	"github.com/AnatoleLucet/loom-term/core/gfx"

	"github.com/AnatoleLucet/go-opentui"
)

type RootElement struct {
	*BaseElement

	initialized bool

	cb  *gfx.CommandBuffer
	rdr *opentui.Renderer
}

func NewRootElement(typ RenderType) (e *RootElement, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Root: %w: %w", ErrFailedToInitializeRoot, err)
		}
	}()

	base, err := NewBaseElement()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	e = &RootElement{
		BaseElement: base,
		cb:          gfx.NewCommandBuffer(ctx),
		rdr:         opentui.NewRenderer(1, 1),
	}

	rc, err := NewRenderContext(ctx, typ, e)
	if err != nil {
		cancel()
		return nil, err
	}
	e.setContextUnsafe(rc)

	go e.listenToMouseEvents(ctx)
	go e.listenToKeyboardEvents(ctx)
	go e.listenToResizeEvents(ctx)
	go e.listenToCapabilites(ctx)
	go e.listenToExitEvents(ctx)

	e.OnDestroy(func() {
		cancel()
		e.rdr.Close()
	})

	return e, nil
}

func (r *RootElement) RenderContext() *RenderContext {
	return r.rdrctx
}
