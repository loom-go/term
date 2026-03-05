package elements

import (
	"context"
	"fmt"
	"runtime"
	"sync"
)

const maxNestedSchedule = 100

type update struct {
	fn   func() error
	done chan struct{}
}

type SchedulerOptions struct {
	Context context.Context
	Errors  chan error

	Render     func() error
	PostRender func() error
}

// Scheduler is a simple orchestrator for updates, renders, and accesses of the element tree.
// It runs accesses synchronously, blocking if a flush (updates+render) is in progress.
// Updates are scheduled async and executed in fifo during the flush. Then the flush triggers a render at the end.
// If an update is scheduled during the render, it will be deffered to the next flush, ensuring the tree can we renderer
// without snapshotting or having a global lock.
type Scheduler struct {
	accessMu sync.RWMutex
	queueMu  sync.Mutex
	ctx      context.Context
	errors   chan error

	render     func() error
	postRender func() error

	queue []*update

	flushing   bool
	batchDepth int
}

func NewScheduler(opts SchedulerOptions) *Scheduler {
	return &Scheduler{
		ctx:        opts.Context,
		errors:     opts.Errors,
		render:     opts.Render,
		postRender: opts.PostRender,
	}
}

func (s *Scheduler) Access(fn func()) {
	s.accessMu.RLock()
	defer s.accessMu.RUnlock()
	fn()
}

func (s *Scheduler) Update(fn func() error) {
	s.enqueue(fn)
	s.Flush()
}

func (s *Scheduler) UpdateSync(fn func() error) {
	done := s.enqueue(fn)
	s.Flush()
	<-done
}

func (s *Scheduler) PushHold() {
	s.queueMu.Lock()
	s.batchDepth++
	s.queueMu.Unlock()
}

func (s *Scheduler) PopHold() {
	s.queueMu.Lock()
	s.batchDepth--
	s.queueMu.Unlock()

	s.Flush()
}

func (s *Scheduler) Batch(fn func()) {
	s.PushHold()
	fn()
	s.PopHold()
}

func (s *Scheduler) Flush() {
	s.queueMu.Lock()
	if s.batchDepth > 0 || s.flushing {
		s.queueMu.Unlock()
		return
	}
	s.flushing = true
	s.queueMu.Unlock()

	// block accessMu to prevent element reads during flush
	// do it synchronously to avoid the gap between flush and the goroutine execution
	s.accessMu.Lock()
	go s.flush()
}

func (s *Scheduler) enqueue(fn func() error) chan struct{} {
	u := &update{fn: fn, done: make(chan struct{}, 1)}

	s.queueMu.Lock()
	s.queue = append(s.queue, u)
	s.queueMu.Unlock()

	return u.done
}

func (s *Scheduler) flush() {
	defer func() {
		s.queueMu.Lock()
		s.flushing = false
		queue := s.queue
		s.queueMu.Unlock()

		s.accessMu.Unlock()

		err := s.postRender()
		if err != nil {
			s.emitError(err)
		}

		// keep going if more updates came in during render
		if len(queue) > 0 {
			runtime.Gosched() // make sure other goroutines have time to enqueue
			s.Flush()
		}
	}()

	for i := range maxNestedSchedule {
		s.queueMu.Lock()
		queue := s.queue
		s.queue = nil
		s.queueMu.Unlock()

		if len(queue) == 0 {
			if i == 0 {
				return // early return if we somehow started a flush with zero update to run
			}

			break
		}

		s.drain(queue)
	}

	err := s.render()
	if err != nil {
		s.emitError(err)
	}
}

func (s *Scheduler) drain(queue []*update) {
	for _, u := range queue {
		if s.ctx.Err() != nil {
			u.done <- struct{}{}
			continue
		}

		s.run(u)
	}
}

func (s *Scheduler) run(u *update) {
	defer func() {
		if r := recover(); r != nil {
			s.emitError(fmt.Errorf("%w: %v", ErrPanicDuringUpdate, r))
		}

		u.done <- struct{}{}
	}()

	err := u.fn()
	if err != nil {
		s.emitError(err)
	}
}

func (s *Scheduler) emitError(err error) {
	select {
	case s.errors <- err:
	default:
	}
}
