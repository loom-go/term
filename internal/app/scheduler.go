package app

import (
	"fmt"
	"sync"

	"github.com/AnatoleLucet/loom-term/core/types"
	"github.com/AnatoleLucet/loom/signals"
)

// Scheduler is responsible for scheduling renders when elements are updated without over or under rendering.
type Scheduler struct {
	mu sync.RWMutex

	ctx types.RenderContext

	render   func() error
	deffered []func() error

	scheduled bool
	clock     int
	holdDepth int

	errorChan chan error
}

func NewScheduler(ctx types.RenderContext, render func() error) *Scheduler {
	return &Scheduler{
		ctx:       ctx,
		render:    render,
		errorChan: make(chan error, 1),
	}
}

func (s *Scheduler) Time() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clock
}

func (s *Scheduler) Errors() <-chan error {
	return s.errorChan
}

func (s *Scheduler) PushHold() {
	s.mu.Lock()
	s.holdDepth++
	s.mu.Unlock()
}

func (s *Scheduler) PopHold() error {
	s.mu.Lock()
	if s.holdDepth > 0 {
		s.holdDepth--
	}
	s.mu.Unlock()

	return s.Schedule()
}

func (s *Scheduler) Defer(f func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.deffered = append(s.deffered, f)
}

func (s *Scheduler) Schedule() error {
	s.mu.Lock()
	if s.holdDepth > 0 || s.scheduled {
		s.mu.Unlock()
		return nil
	}

	s.scheduled = true
	s.mu.Unlock()

	// enqueue directly for the first render because the reactive system is not being updated yet,
	// we're just initializing the tree
	if s.clock == 0 {
		go s.doRender()
		return nil
	}

	// else, it means a signal has been updated and the reactive system is flusing.
	// so wait for every render effects (element updates) to settle to enqueue only one render afterwards
	signals.OnRenderSettled(func() {
		go s.doRender()
	})

	return nil
}

func (s *Scheduler) doRender() {
	s.ctx.LockRender()
	defer s.ctx.UnlockRender()

	for s.scheduled {
		err := s.render()

		s.mu.Lock()
		s.clock++
		s.scheduled = false
		s.mu.Unlock()

		err = s.drainDeferred() // might reschedule an update

		if err != nil {
			s.errorChan <- fmt.Errorf("Scheduler: %w", err)
			return
		}
	}
}

func (s *Scheduler) drainDeferred() error {
	s.mu.Lock()
	deffered := s.deffered
	s.deffered = nil
	s.mu.Unlock()

	for _, f := range deffered {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}
