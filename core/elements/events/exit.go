package events

import (
	"context"
	"github.com/AnatoleLucet/loom-term/core/sync"
	"os"
	"os/signal"
	"syscall"
)

type ExitListener struct {
	ctx    context.Context
	events *sync.Broadcaster[*EventExit]
}

func NewExitListener(ctx context.Context) *ExitListener {
	listener := &ExitListener{
		ctx:    ctx,
		events: sync.NewBroadcaster[*EventExit](ctx),
	}

	go listener.watch()

	return listener
}

func (l *ExitListener) Listen(ctx context.Context) <-chan *EventExit {
	return l.events.Listen(ctx)
}

func (l *ExitListener) watch() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	for {
		select {
		case <-l.ctx.Done():
			return
		case sig := <-ch:
			s, ok := l.toExitSignal(sig)
			if !ok {
				continue
			}

			l.events.Broadcast(&EventExit{Signal: s})
		}
	}
}

func (l *ExitListener) toExitSignal(sig os.Signal) (ExitSignal, bool) {
	switch sig {
	case syscall.SIGINT:
		return ExitSigInt, true
	case syscall.SIGTERM:
		return ExitSigTerm, true
	case syscall.SIGQUIT:
		return ExitSigQuit, true
	case syscall.SIGHUP:
		return ExitSigHup, true
	default:
		return 0, false
	}
}
