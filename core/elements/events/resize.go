package events

import (
	"context"
	"github.com/AnatoleLucet/loom-term/core/sync"
	"github.com/AnatoleLucet/loom-term/core/term"
	"os"
	"os/signal"
	"syscall"
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

func (l *ResizeListener) watch() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ch:
			width, height, err := term.Size()
			if err == nil {
				l.events.Broadcast(&EventResize{Width: width, Height: height})
			}
		}
	}
}
