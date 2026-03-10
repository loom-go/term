//go:build windows

package events

import (
	"time"

	"github.com/loom-go/term/core/term"
)

func (l *ResizeListener) watch() {
	lastWidth, lastHeight, err := term.Size()
	if err != nil {
		lastWidth, lastHeight = 0, 0
	}

	for range time.Tick(time.Millisecond * 100) {
		select {
		case <-l.ctx.Done():
			return
		default:
		}

		width, height, err := term.Size()
		if err != nil {
			continue
		}

		if width != lastWidth || height != lastHeight {
			lastWidth, lastHeight = width, height
			l.events.Broadcast(&EventResize{Width: width, Height: height})
		}
	}
}
