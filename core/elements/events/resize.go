package events

import (
	"context"
	"github.com/loom-go/term/core/sync"
)

type ResizeListener struct {
	ctx    context.Context
	events *sync.Broadcaster[*EventResize]
}

func NewResizeListener(ctx context.Context) *ResizeListener {
	listener := &ResizeListener{
		ctx:    ctx,
		events: sync.NewBroadcaster[*EventResize](ctx),
	}

	go listener.watch()

	return listener
}

func (l *ResizeListener) Listen(ctx context.Context) <-chan *EventResize {
	return l.events.Listen(ctx)
}
