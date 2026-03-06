package term

import (
	"context"
	"fmt"
	"sync"

	"github.com/AnatoleLucet/loom"
	appctx "github.com/AnatoleLucet/loom-term/components/context"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom/signals"
)

type RenderType = core.RenderType

const RenderInline RenderType = core.RenderInline
const RenderFullscreen RenderType = core.RenderFullscreen

type AppContext = appctx.Context

func Context() *AppContext {
	ctx, err := appctx.Get()
	if err != nil {
		// Context should only be used in a reactive scope.
		// If it's not, we have good reasons to *panic*
		panic(fmt.Errorf("term.Context: %w", err))
	}

	return ctx
}

type App struct {
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	rdr   *core.Renderer
	owner *signals.Owner

	running bool
	errors  chan any
}

func NewApp() *App {
	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		ctx:     ctx,
		cancel:  cancel,
		running: false,
		errors:  make(chan any, 1),
	}
}

func (a *App) Run(typ RenderType, fn func() loom.Node) <-chan any {
	a.mu.Lock()
	defer a.mu.Unlock()

	var err error
	defer func() {
		if err != nil {
			a.errors <- fmt.Errorf("App: %w", err)
			a.close()
		}
	}()

	if a.running {
		err = ErrAppAlreadyRunning
		return a.Errors()
	}

	a.rdr, err = core.NewRenderer(typ)
	if err != nil {
		return a.Errors()
	}

	err = a.render(fn)
	if err != nil {
		return a.Errors()
	}

	// syncronously cancel the ctx if the renderer closed itself
	// so the user can handle it with app.Close()
	a.rdr.OnClose(a.cancel)

	// forward errors from the renderer
	go forward(a.ctx, a.rdr.Errors(), a.errors)

	a.running = true
	return a.Errors()
}

func (a *App) Errors() <-chan any {
	out := make(chan any)

	go func() {
		defer close(out)

		for {
			select {
			case <-a.ctx.Done():
				return
			case err, ok := <-a.errors:
				if !ok {
					return
				}

				select {
				case out <- err:
				case <-a.ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

func (a *App) render(fn func() loom.Node) error {
	appctx := appctx.New(a.rdr)

	root, err := newRootNode(appctx, fn)
	if err != nil {
		return fmt.Errorf("failed to create root node:  %w", err)
	}

	a.owner, err = loom.Render(nil, root)
	if err != nil {
		return fmt.Errorf("failed on initial render: %w", err)
	}

	// handle panics in the reactive system
	a.owner.OnError(a.emitError)

	return nil
}

func (a *App) Stop() {
	a.Close()
}

func (a *App) Close() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return
	}
	a.running = false

	a.close()
}

func (a *App) close() {
	a.cancel()

	if a.rdr != nil {
		a.rdr.Close()
		a.rdr = nil
	}

	if a.owner != nil {
		a.owner.Dispose()
		a.owner = nil
	}

	close(a.errors)
}

func (a *App) emitError(err any) {
	if !a.running || a.ctx.Err() != nil {
		return
	}

	select {
	case a.errors <- err:
	case <-a.ctx.Done():
	}
}

func forward[t any](ctx context.Context, source <-chan t, target chan<- any) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-source:
			select {
			case target <- v:
			case <-ctx.Done():
				return
			}
		}
	}
}
