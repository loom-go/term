package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/AnatoleLucet/loom-term/core/elements"
	"github.com/AnatoleLucet/loom-term/core/events"
	"github.com/AnatoleLucet/loom-term/core/render"
	"github.com/AnatoleLucet/loom-term/core/stdio"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type TermRuntime struct {
	mu sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc

	renderCtx types.RenderContext

	root     types.Element
	renderer types.Renderer
	events   types.EventListener

	errorsCh chan error

	closed bool
}

func NewRuntime(typ types.RenderType) (rt *TermRuntime, err error) {
	rt = &TermRuntime{}
	rt.ctx, rt.cancel = context.WithCancel(context.Background())

	rt.renderCtx = NewRenderContext(rt)

	rt.events = events.NewListener()

	rt.root, err = elements.NewElement(rt.renderCtx)
	if err != nil {
		return nil, fmt.Errorf("Runtime: %w", err)
	}

	rt.renderer, err = render.NewRenderer(typ, rt.renderCtx)
	if err != nil {
		return nil, fmt.Errorf("Runtime: %w", err)
	}

	go rt.watchCapabilites(rt.ctx)
	go rt.watchResize(rt.ctx)
	go rt.watchMouse(rt.ctx)
	// go rt.watchExit(rt.ctx)

	return rt, nil
}

func (a *TermRuntime) Root() types.Element {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.root
}

func (a *TermRuntime) RenderContext() types.RenderContext {
	return a.renderCtx
}

func (a *TermRuntime) Render() error {
	a.renderCtx.LockRender()
	defer a.renderCtx.UnlockRender()

	return a.RenderUnsafe()
}

func (a *TermRuntime) RenderUnsafe() error {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if err := a.guard(); err != nil {
		return err
	}

	err := a.renderer.Render(a.root)
	if err != nil {
		return fmt.Errorf("Runtime: %w: %w", types.ErrFailedToRenderFrame, err)
	}

	return nil
}

func (a *TermRuntime) Errors() <-chan error {
	return a.errorsCh
}

func (a *TermRuntime) Events() (events <-chan types.Event, stop func()) {
	ctx, stop := context.WithCancel(context.Background())
	events, _ = a.events.Listen(ctx)

	return events, stop
}

func (a *TermRuntime) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.closed {
		return nil
	}
	a.closed = true

	a.renderCtx.LockRender()
	a.cancel()
	a.renderer.Close()
	a.root.Destroy()
	a.renderCtx.UnlockRender()

	return nil
}

func (a *TermRuntime) IsClosed() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.closed
}

func (a *TermRuntime) guard() error {
	if a.closed {
		return types.ErrUsingClosedRuntime
	}
	return nil
}

func (a *TermRuntime) watchCapabilites(ctx context.Context) {
	stdin := stdio.Stdin.Listen(1024)
	consumer := stdio.NewBufferedConsumer(func(buf []byte) (int, bool) {
		err := a.renderer.UpdateCapabilities(buf)
		if err != nil {
			return 0, false // wait for more data
		}

		return len(buf), true
	})

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-stdin:
			consumer.Feed(data)
		}
	}
}

func (a *TermRuntime) watchResize(ctx context.Context) {
	eventsCh, _ := a.events.Listen(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-eventsCh:
			switch e := event.(type) {
			case *types.EventResize:
				a.root.SetWidth(e.Width)

				// todo: we still have some issues when resizing the terminal in inline mode
				if a.renderer.Type() == types.RenderTypeFullscreen {
					a.root.SetHeight(e.Height)
				}

				err := a.Render()
				if err != nil {
					a.errorsCh <- err
				}
			}
		}
	}
}

func (a *TermRuntime) watchMouse(ctx context.Context) {
	eventsCh, _ := a.events.Listen(ctx)

	var current types.Element
	var captured types.Element

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-eventsCh:
			switch e := event.(type) {
			case *types.EventMouse:
				e = &types.EventMouse{
					X:      e.X,
					Y:      e.Y - a.renderer.RenderOffset(), // normalize mouse coordinates before propagating to elements
					Ctrl:   e.Ctrl,
					Alt:    e.Alt,
					Shift:  e.Shift,
					Action: e.Action,
					Button: e.Button,
				}

				elem := a.checkHit(e.X, e.Y)
				if elem == nil {
					if current != nil {
						current.HandleMouseLeave(e)
						current = nil
					}

					continue
				}

				if e.Action == types.MouseActionMove {
					if current != nil && current == elem {
						elem.HandleMouseMove(e)
						current = elem
					} else {
						if current != nil && current != elem {
							current.HandleMouseLeave(e)
						}

						elem.HandleMouseEnter(e)
						current = elem
					}
				}

				if e.Action == types.MouseActionPress {
					elem.HandleMousePress(e)
					captured = elem
				}

				if e.Action == types.MouseActionRelease {
					elem.HandleMouseRelease(e)
					captured = nil
				}

				if e.Action == types.MouseActionScroll {
					elem.HandleMouseScroll(e)
				}

				if e.Action == types.MouseActionDrag {
					if captured == nil {
						elem.HandleMouseDrag(e)
					} else {
						captured.HandleMouseDrag(e)
					}
				}
			}
		}
	}
}

func (a *TermRuntime) checkHit(x, y int) types.Element {
	return a.renderer.CheckHit(x, y)
}

func (a *TermRuntime) addToHitGrid(element types.Element, x, y, width, height int) {
	a.renderer.AddToHitGrid(element, x, y, width, height)
}

func (a *TermRuntime) pushHitGridScissors(x, y, width, height int) {
	a.renderer.PushHitGridScissors(x, y, width, height)
}

func (a *TermRuntime) popHitGridScissors() {
	a.renderer.PopHitGridScissors()
}
