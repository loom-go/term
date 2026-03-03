package app

import (
	"context"

	"github.com/AnatoleLucet/loom-term/core"
	"github.com/AnatoleLucet/loom/signals"
)

var Context = signals.NewContext[*AppContext](nil)

func GetContext() (*AppContext, error) {
	ctx := Context.Value()
	if ctx == nil {
		return nil, ErrAppNotInitialized
	}

	return ctx, nil
}

type AppContext struct {
	rdr       *core.Renderer
	scheduler *Scheduler
}

func NewAppContext(ctx context.Context, rdr *core.Renderer) *AppContext {
	return &AppContext{
		rdr:       rdr,
		scheduler: NewScheduler(ctx, rdr.ScheduleRender),
	}
}

func (ac *AppContext) Root() core.Element {
	return ac.rdr.Root()
}

func (ac *AppContext) PushRenderHold() {
	ac.scheduler.PushHold()
}

func (ac *AppContext) PopRenderHold() {
	ac.scheduler.PopHold()
}

// ScheduleRender schedules a render to be done when the holds are released,
// and orchetrated with the reactive system (once render effects have been executed)
func (ac *AppContext) RequestRender() {
	ac.scheduler.ScheduleRender()
}

// ScheduleRender is a direct call to the renderer to render immediately, bypassing holds and the reactive system.
func (ac *AppContext) ScheduleRender() {
	ac.rdr.ScheduleRender()
}
