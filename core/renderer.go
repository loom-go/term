package core

import (
	"context"
	"github.com/AnatoleLucet/loom-term/core/elements"
)

type RenderType = elements.RenderType

const (
	RenderInline     RenderType = elements.RenderTypeInline
	RenderFullscreen RenderType = elements.RenderTypeFullscreen
)

type RootElement = *elements.RootElement
type RenderContext = *elements.RenderContext

type Renderer struct {
	root *elements.RootElement

	ctx context.Context
}

func NewRenderer(typ RenderType) (*Renderer, error) {
	root, err := elements.NewRootElement(typ)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	r := &Renderer{root: root, ctx: ctx}

	r.Root().OnDestroy(cancel)

	return r, nil
}

func (r *Renderer) Root() RootElement {
	return r.root
}

func (r *Renderer) Render() error {
	return r.root.RenderContext().Render()
}

func (r *Renderer) ScheduleRender() {
	r.root.RenderContext().ScheduleRender()
}

func (r *Renderer) Errors() <-chan error {
	return r.root.RenderContext().Errors()
}

func (r *Renderer) Wait() {
	<-r.ctx.Done()
}

func (r *Renderer) Close() {
	r.root.Destroy()
}

func (r *Renderer) Closed() <-chan struct{} {
	return r.ctx.Done()
}
