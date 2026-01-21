package debug

import (
	"context"
	"sync"
)

type Emitter[T any] struct {
	mu        sync.RWMutex
	listeners map[chan T]struct{}
}

func NewEmitter[T any]() *Emitter[T] {
	return &Emitter[T]{
		listeners: make(map[chan T]struct{}),
	}
}

func (e *Emitter[T]) Subscribe(ctx context.Context, buffer int) <-chan T {
	ch := make(chan T, buffer)

	e.mu.Lock()
	e.listeners[ch] = struct{}{}
	e.mu.Unlock()

	go func() {
		<-ctx.Done()
		e.mu.Lock()
		delete(e.listeners, ch)
		close(ch)
		e.mu.Unlock()
	}()

	return ch
}

func (e *Emitter[T]) Emit(v T) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for ch := range e.listeners {
		select {
		case ch <- v:
		default:
			// drop oldest and send new value
			<-ch
			ch <- v
		}
	}
}
