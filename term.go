package term

import (
	"context"
	"fmt"
	"sync"

	"github.com/AnatoleLucet/loom"
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/app"
	"github.com/AnatoleLucet/loom/signals"
)

type RenderType = core.RenderType

const RenderInline RenderType = core.RenderInline
const RenderFullscreen RenderType = core.RenderFullscreen

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
		return a.errors
	}

	a.rdr, err = core.NewRenderer(typ)
	if err != nil {
		return a.errors
	}

	err = a.render(fn)
	if err != nil {
		return a.errors
	}

	// syncronously destroy the app if the root destroys itself
	a.rdr.Root().OnDestroy(a.Close)

	// forward errors from the renderer
	go forward(a.ctx, a.rdr.Errors(), a.errors)

	a.running = true
	return a.errors
}

func (a *App) render(fn func() loom.Node) error {
	appctx := app.NewAppContext(a.ctx, a.rdr)

	root, err := newRootNode(a.ctx, appctx, fn)
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
	// make sure to cancel the ctx BEFORE unmounting the tree,
	// else some goroutine might squeeze in an update/render on the destroyed tree before the ctx is cancelled.
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
