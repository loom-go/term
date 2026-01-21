package types

import "context"

type Event interface {
	ShouldPropagate() bool
	StopPropagation()
}

type BaseEvent struct {
	stopped bool
}

func (e *BaseEvent) ShouldPropagate() bool {
	return !e.stopped
}

func (e *BaseEvent) StopPropagation() {
	e.stopped = false
}

type EventListener interface {
	Listen(ctx context.Context) (<-chan Event, <-chan error)
	Close()
}
