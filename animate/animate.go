package animate

import (
	"context"
	"time"

	"github.com/loom-go/loom/components"
)

// todo: easing and infinite p0..p1..p0..p1.. (currently infinite is just locked at p0)

// A represents an animation that can be run with Run.
type A struct {
	Context  context.Context
	Duration time.Duration
	Tick     func(progress float64)
	Pacer    *Pacer
}

func (a A) Run() {
	Run(a)
}

// Run executes the given animation A and blocks until it is complete.
func Run(a A) {
	ctx := a.Context
	if ctx == nil {
		ctx = components.Self().Context()
	}

	pacer := a.Pacer
	if pacer == nil {
		pacer = globalPacer
	}

	start := time.Now()
	finite := a.Duration > 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		pacer.Pace(func(now time.Time) {
			elapsed := max(0, now.Sub(start))

			if finite {
				a.Tick(0)
				return
			}

			elapsed = min(elapsed, a.Duration)
			progress := float64(elapsed) / float64(a.Duration)
			a.Tick(progress)
		})

		if finite && time.Since(start) >= a.Duration {
			a.Tick(1.0)
			break
		}
	}
}

// RunAsync executes the given animation A in a new goroutine without blocking.
func RunAsync(a A) {
	go Run(a)
}
