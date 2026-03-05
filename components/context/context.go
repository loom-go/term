package appctx

import (
	"context"

	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom/signals"
)

var AppContext = signals.NewContext[*Context](nil)

func Get() (*Context, error) {
	ctx := AppContext.Value()
	if ctx == nil {
		return nil, ErrAppNotInitialized
	}

	return ctx, nil
}

type Context struct {
	rdr *core.Renderer
}

func New(ctx context.Context, rdr *core.Renderer) *Context {
	return &Context{rdr: rdr}
}

func (ac *Context) Root() core.Element {
	return ac.rdr.Root()
}

func (ac *Context) RenderType() core.RenderType {
	return ac.rdr.RenderType()
}

func (ac *Context) BatchRender(fn func() error) (err error) {
	ac.rdr.Batch(func() { err = fn() })

	return
}
