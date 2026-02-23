package sync

import (
	"context"
	"sync"
)

type Broadcaster[T any] struct {
	mu  sync.RWMutex
	ctx context.Context

	listeners []chan T
}

func NewBroadcaster[T any](ctx context.Context) *Broadcaster[T] {
	return &Broadcaster[T]{ctx: ctx}
}

func (m *Broadcaster[T]) Listen(ctx context.Context) <-chan T {
	ch := make(chan T)

	m.mu.Lock()
	m.listeners = append(m.listeners, ch)
	m.mu.Unlock()

	go func() {
		<-ctx.Done()
		m.mu.Lock()
		defer m.mu.Unlock()

		for i, listener := range m.listeners {
			if listener == ch {
				m.listeners = append(m.listeners[:i], m.listeners[i+1:]...)
				break
			}
		}
	}()

	return ch
}

func (m *Broadcaster[T]) Broadcast(event T) {
	select {
	case <-m.ctx.Done():
		return
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, listener := range m.listeners {
		select {
		case listener <- event:
		default:
		}
	}
}
