package app

import (
	"context"
	"sync"

	"github.com/AnatoleLucet/loom/signals"
)

// Scheduler is responsible for scheduling renders when elements are updated without over or under rendering.
type Scheduler struct {
	mu sync.RWMutex

	ctx context.Context

	// render   func() error
	render func()

	deffered []func() error

	scheduled bool
	clock     int
	holdDepth int

	errors chan error
}

func NewScheduler(ctx context.Context, render func()) *Scheduler {
	return &Scheduler{
		ctx:    ctx,
		render: render,
		errors: make(chan error, 1),
	}
}

func (s *Scheduler) Time() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clock
}

func (s *Scheduler) Errors() <-chan error {
	return s.errors
}

func (s *Scheduler) PushHold() {
	s.mu.Lock()
	s.holdDepth++
	s.mu.Unlock()
}

func (s *Scheduler) PopHold() {
	s.mu.Lock()
	if s.holdDepth > 0 {
		s.holdDepth--
	}
	s.mu.Unlock()

	s.Schedule()
}

func (s *Scheduler) Defer(f func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deffered = append(s.deffered, f)
}

func (s *Scheduler) Schedule() {
	s.mu.Lock()
	if s.holdDepth > 0 || s.scheduled {
		s.mu.Unlock()
		return
	}
	s.scheduled = true
	s.mu.Unlock()

	// enqueue directly for the first render because the reactive system is not being updated yet,
	// we're just initializing the tree
	if s.clock == 0 {
		s.doRender()
	}

	// else, it means a signal has been updated and the reactive system is flusing.
	// so wait for every render effects (element updates) to settle to enqueue only one render afterwards
	signals.OnRenderSettled(func() {
		s.doRender()
	})

}

func (s *Scheduler) doRender() {
	// for s.scheduled {
	// 	select {
	// 	case <-s.ctx.Done():
	// 		return
	// 	default:
	// 	}

	s.mu.Lock()
	s.clock++
	s.scheduled = false
	s.mu.Unlock()

	s.render()

	// 	if s.scheduled {
	// 		core.LogDebugf("Render scheduled during render, deferring next render")
	// 	}
	// }
}

func (s *Scheduler) drainDeferred() error {
	s.mu.Lock()
	deffered := s.deffered
	s.deffered = nil
	s.mu.Unlock()

	for _, f := range deffered {
		select {
		case <-s.ctx.Done():
			return nil
		default:
		}

		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
