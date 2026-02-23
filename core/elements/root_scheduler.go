package elements

import (
	"context"
	"fmt"
	"sync"
)

type taskType int

const (
	taskUpdate taskType = iota
	taskRender
)

type task struct {
	typ    taskType
	fn     func() error
	result chan error
}

// Scheduler is a simpler single-threaded fifo loop that can either run an update or a render.
// If an update is scheduled mid-render, it will be deffered to after the render is done.
// This ensure we can have unblocking updates during rendering, while not having to snapshot or lock the whole tree for each render.
type Scheduler struct {
	mu  sync.RWMutex
	ctx context.Context

	task  chan *task // Internal coordination only
	queue []*task    // The actual ordered queue

	rendering bool
	pending   bool // Signal that queue has work
}

func NewScheduler(ctx context.Context) *Scheduler {
	s := &Scheduler{
		ctx:  ctx,
		task: make(chan *task, 1024),
	}

	go s.loop()

	return s
}

func (s *Scheduler) Schedule(typ taskType, fn func() error) <-chan error {
	result := make(chan error, 1)
	task := &task{typ: typ, fn: fn, result: result}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.queue = append(s.queue, task)
	if s.rendering && typ == taskUpdate {
		return result
	}

	if !s.pending {
		s.pending = true
		select {
		case s.task <- task: // wake up
		default:
		}
	}

	return result
}

func (s *Scheduler) loop() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.task:
			s.processQueue()
		}
	}
}

func (s *Scheduler) processQueue() {
	for {
		s.mu.Lock()
		if len(s.queue) == 0 {
			s.pending = false
			s.mu.Unlock()
			return
		}

		task := s.queue[0]
		s.queue = s.queue[1:]

		isRender := task.typ == taskRender
		if isRender {
			s.rendering = true
		}
		s.mu.Unlock()

		s.run(task)

		if isRender {
			s.mu.Lock()
			s.rendering = false
			s.mu.Unlock()
		}
	}
}

func (s *Scheduler) send(task *task) {
	select {
	case s.task <- task:
	case <-s.ctx.Done():
		task.result <- nil
	}
}

func (s *Scheduler) run(task *task) {
	defer func() {
		// just in case since we're often doing CGO calls that could panic somehow
		if r := recover(); r != nil {
			task.result <- fmt.Errorf("%v", r)
		}
	}()

	task.result <- task.fn()
}

func (s *Scheduler) enqueue(task *task) {
	s.queue = append(s.queue, task)
}

func (s *Scheduler) flush(queue []*task) {
	for _, task := range queue {
		if err := s.ctx.Err(); err != nil {
			task.result <- err
			continue
		}

		s.run(task)
	}
}
