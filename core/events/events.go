package events

import (
	"context"
	"os"
	"os/signal"
	"slices"
	"sync"
	"syscall"

	"github.com/AnatoleLucet/loom-term/core/stdio"
	"github.com/AnatoleLucet/loom-term/core/terminal"
	"github.com/AnatoleLucet/loom-term/core/types"
)

type listenner struct {
	events chan types.Event
	errors chan error
}

type Listener struct {
	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	events chan types.Event
	errors chan error

	listenners []listenner
}

func NewListener() *Listener {
	l := &Listener{}

	l.ctx, l.cancel = context.WithCancel(context.Background())

	go l.listenStdin(l.ctx)
	go l.listenResize(l.ctx)
	go l.listenExit(l.ctx)

	return l
}

func (l *Listener) Listen(ctx context.Context) (<-chan types.Event, <-chan error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	events := make(chan types.Event, 100)
	errors := make(chan error, 100)

	l.listenners = append(l.listenners, listenner{
		events: events,
		errors: errors,
	})

	go func() {
		<-ctx.Done()
		l.mu.Lock()
		defer l.mu.Unlock()

		l.listenners = slices.DeleteFunc(l.listenners, func(listenner listenner) bool {
			return listenner.events == events && listenner.errors == errors
		})

		close(events)
		close(errors)
	}()

	return events, errors
}

func (l *Listener) Close() {
	l.cancel()
}

func (l *Listener) emitEvent(event types.Event) {
	if event == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for _, listenner := range l.listenners {
		select {
		case listenner.events <- event:
		default:
			// drop oldest and send new value
			<-listenner.events
			listenner.events <- event
		}
	}
}

func (l *Listener) emitError(err error) {
	if err == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for _, listenner := range l.listenners {
		select {
		case listenner.errors <- err:
		default:
			// drop oldest and send new value
			<-listenner.errors
			listenner.errors <- err
		}
	}
}

func (l *Listener) listenStdin(ctx context.Context) {
	stdin := stdio.Stdin.Listen(256)

	consumer := stdio.NewBufferedConsumer(func(buf []byte) (consumed int, complete bool) {
		event, consumed, ok := l.processStdin(buf)
		if !ok {
			return 0, false
		}

		l.emitEvent(event)
		return consumed, true
	})

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-stdin:
			if !ok {
				return
			}
			consumer.Feed(data)
		}
	}
}

func (l *Listener) listenResize(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGWINCH)

	for {
		select {
		case <-ctx.Done():
			return
		case <-sigChan:
			width, height, err := terminal.Size()
			if err != nil {
				l.emitError(err)
				continue
			}
			l.emitEvent(&types.EventResize{Width: width, Height: height})
		}
	}
}

func (l *Listener) listenExit(ctx context.Context) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	for {
		select {
		case <-ctx.Done():
			return
		case sig := <-sigChan:
			s, ok := toExitSignal(sig)
			if !ok {
				continue
			}
			l.emitEvent(&types.EventExit{Signal: s})
		}
	}
}

func (l *Listener) processStdin(buf []byte) (event types.Event, consumed int, ok bool) {
	if mouseEvent, mouseConsumed, mouseOk := l.processMouseEvent(buf); mouseOk {
		return mouseEvent, mouseConsumed, true
	}

	if pasteEvent, pasteConsumed, pasteOk := l.processPasteEvent(buf); pasteOk {
		return pasteEvent, pasteConsumed, true
	}

	if keyEvent, keyConsumed, keyOk := l.processKeyEvent(buf); keyOk {
		return keyEvent, keyConsumed, true
	}

	return nil, 0, false
}

func (l *Listener) processMouseEvent(buf []byte) (types.Event, int, bool) {
	const mouseEventX10Len = 6
	if len(buf) >= mouseEventX10Len && buf[0] == '\x1b' && buf[1] == '[' {
		switch buf[2] {
		case 'M':
			return parseX10MouseEvent(buf), mouseEventX10Len, true
		case '<':
			// need at least: ESC[<0;0;0M = 10 bytes
			if len(buf) < 10 {
				return nil, 0, false
			}

			for i := 4; i < len(buf); i++ {
				if buf[i] == 'M' || buf[i] == 'm' {
					return parseSGRMouseEvent(buf[:i+1]), i + 1, true
				}
			}

			return nil, 0, false
		}
	}

	return nil, 0, false
}

func (l *Listener) processPasteEvent(buf []byte) (types.Event, int, bool) {
	if len(buf) < 6 || string(buf[:6]) != "\x1b[200~" {
		return nil, 0, false
	}

	paste, consumed, ok := parseBracketedPaste(buf)
	if !ok {
		return nil, 0, false
	}
	if paste == nil {
		return nil, consumed, true
	}

	return paste, consumed, true
}

func (l *Listener) processKeyEvent(buf []byte) (types.Event, int, bool) {
	key, consumed, ok := parseKeyEvent(buf)
	if !ok {
		return nil, 0, false
	}
	if key == nil {
		return nil, consumed, true
	}

	return key, consumed, true
}
