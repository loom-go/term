//go:build !windows

package events

import (
	"github.com/loom-go/term/core/term"
	"os"
	"os/signal"
	"syscall"
)

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
