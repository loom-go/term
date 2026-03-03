package app

import (
	"context"
	"sync"

	"github.com/AnatoleLucet/loom/signals"
)

// Scheduler is responsible for scheduling renders within the reactive system
// when elements are updated without over or under rendering
type Scheduler struct {
	mu sync.RWMutex

	ctx context.Context

	render func()

	scheduled bool
	clock     int
	holdDepth int
}

func NewScheduler(ctx context.Context, render func()) *Scheduler {
	return &Scheduler{
		ctx:    ctx,
		render: render,
	}
}

func (s *Scheduler) Time() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clock
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

	s.ScheduleRender()
}

func (s *Scheduler) ScheduleRender() {
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
	s.mu.Lock()
	s.clock++
	s.scheduled = false
	s.mu.Unlock()

	s.render()
}
