package events

import (
	"context"
	"github.com/AnatoleLucet/loom-term/core/stdio"
	"github.com/AnatoleLucet/loom-term/core/sync"
	"github.com/AnatoleLucet/loom-term/core/term"
)

type CapabilitiesListener struct {
	ctx    context.Context
	events *sync.Broadcaster[*EventCapabilities]
}

func NewCapabilitiesListener(ctx context.Context) *CapabilitiesListener {
	listener := &CapabilitiesListener{
		ctx:    ctx,
		events: sync.NewBroadcaster[*EventCapabilities](ctx),
	}

	go listener.watch()

	return listener
}

func (l *CapabilitiesListener) Listen(ctx context.Context) <-chan *EventCapabilities {
	return l.events.Listen(ctx)
}

func (l *CapabilitiesListener) watch() {
	stdin := stdio.Stdin.Listen(1024)
	events := stdio.NewBufferedConsumer(func(buf []byte) (consumed int, complete bool) {
		event, consumed := l.parseCapabilitiesEvent(buf)
		if event != nil {
			l.events.Broadcast(event)
			return consumed, true
		}

		return consumed, consumed > 0
	})

	for {
		select {
		case <-l.ctx.Done():
			return
		case buf := <-stdin:
			events.Feed(buf)
		}
	}
}

func (l *CapabilitiesListener) parseCapabilitiesEvent(buf []byte) (event *EventCapabilities, consumed int) {
	if term.IsCapabilityResponse(buf) {
		return &EventCapabilities{buf}, len(buf)
	}

	return nil, 1
}
