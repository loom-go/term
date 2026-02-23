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

func (ac *AppContext) Errors() <-chan error {
	return ac.scheduler.Errors()
}

func (ac *AppContext) PushRenderHold() {
	ac.scheduler.PushHold()
}

func (ac *AppContext) PopRenderHold() {
	ac.scheduler.PopHold()
}

func (ac *AppContext) ScheduleRender() {
	ac.scheduler.Schedule()
}
