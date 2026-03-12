package core

import (
	"context"

	"github.com/loom-go/term/core/elements"
	"github.com/loom-go/term/core/term"
)

type RenderType = elements.RenderType

const (
	RenderInline     RenderType = elements.RenderTypeInline
	RenderFullscreen RenderType = elements.RenderTypeFullscreen
)

type RootElement = *elements.RootElement

type Renderer struct {
	ctx    context.Context
	rdrctx *elements.RenderContext
	root   RootElement
	closed bool
}

func NewRenderer(typ RenderType) (*Renderer, error) {
	restore, err := term.Init()
	if err != nil {
		return nil, err
	}

	root, err := elements.NewRootElement(typ)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	rdrctx, err := elements.NewRenderContext(ctx, typ, root)
	if err != nil {
		cancel()
		return nil, err
	}

	root.SetRenderContext(rdrctx)
	root.OnDestroy(func() {
		cancel()
		restore()
	})

	return &Renderer{root: root, rdrctx: rdrctx, ctx: ctx}, nil
}

func (r *Renderer) Root() RootElement      { return r.root }
func (r *Renderer) Errors() <-chan error   { return r.rdrctx.Errors() }
func (r *Renderer) RenderType() RenderType { return r.rdrctx.RenderType() }

func (r *Renderer) PushHold()       { r.rdrctx.PushHold() }
func (r *Renderer) PopHold()        { r.rdrctx.PopHold() }
func (r *Renderer) Batch(fn func()) { r.rdrctx.Batch(fn) }

func (r *Renderer) Closed() <-chan struct{} {
	return r.ctx.Done()
}
func (r *Renderer) OnClose(fn func()) {
	r.root.OnDestroy(fn)
}

func (r *Renderer) Close() {
	if r.closed {
		return
	}
	r.closed = true

	r.root.Destroy()
	r.rdrctx.Settle()
}
