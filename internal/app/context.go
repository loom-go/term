package app

import (
	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom-term/internal/errs"
	"github.com/AnatoleLucet/loom/signals"
)

var Context = signals.NewContext[*AppContext](nil)

func GetContext() (*AppContext, error) {
	ctx := Context.Value()
	if ctx == nil {
		return nil, errs.ErrAppNotInitialized
	}

	return ctx, nil
}

type AppContext struct {
	runtime   core.Runtime
	scheduler *Scheduler
}

func NewAppContext(runtime core.Runtime) *AppContext {
	return &AppContext{
		runtime:   runtime,
		scheduler: NewScheduler(runtime.RenderContext(), runtime.RenderUnsafe),
	}
}

func (ac *AppContext) Root() core.Element {
	return ac.runtime.Root()
}

func (ac *AppContext) Errors() <-chan error {
	return ac.scheduler.Errors()
}

func (ac *AppContext) RenderContext() core.RenderContext {
	return ac.runtime.RenderContext()
}

func (ac *AppContext) PushRenderHold() {
	ac.scheduler.PushHold()
}

func (ac *AppContext) PopRenderHold() error {
	return ac.scheduler.PopHold()
}

func (ac *AppContext) RequestRender() error {
	return ac.scheduler.Schedule()
}

// DoSafely makes shure the given function is executed outside of a render cycle.
// Mainly used for updating elements safely.
func (ac *AppContext) DoSafely(fn func() error) error {
	// try to acquire the render lock to run immediately
	if ac.RenderContext().TryLockRender() {
		defer ac.RenderContext().UnlockRender()
		return fn()
	}

	// else it means we're already rendering, so defer the update to after the current render
	ac.scheduler.Defer(fn)
	return nil
}

func (ac *AppContext) Close() error {
	return ac.runtime.Close()
}
