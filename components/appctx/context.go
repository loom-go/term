package appctx

import (
	components "github.com/loom-go/loom/components"
	"github.com/loom-go/term/core"
)

var _, appContext = components.NewContext[*Context](nil)
var Provider = appContext.Provider

func Get() (*Context, error) {
	ctx := appContext.Get()
	if ctx == nil {
		return nil, ErrAppNotInitialized
	}

	return ctx, nil
}

type Context struct {
	rdr *core.Renderer
}

func New(rdr *core.Renderer) *Context {
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
