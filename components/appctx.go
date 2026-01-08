package components

import (
	termerror "github.com/AnatoleLucet/loom-term/error"
	"github.com/AnatoleLucet/loom-term/internal"
	"github.com/AnatoleLucet/loom/signals"
)

var appctx = signals.NewContext[*AppContext](nil)

func getAppContext() (*AppContext, error) {
	ctx := appctx.Value()
	if ctx == nil {
		return nil, termerror.ErrNoRendererInContext
	}

	return ctx, nil
}

type AppContext struct {
	renderer *internal.TermRenderer
	// debugger *internal.TermDebugger // for fps, logs, etc
}

func NewAppContext(rendererType internal.RendererType) (*AppContext, error) {
	renderer, err := internal.NewRenderer(rendererType)
	if err != nil {
		return nil, err
	}

	return &AppContext{
		renderer: renderer,
	}, nil
}
