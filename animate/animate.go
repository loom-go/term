package animate

import (
	"context"
	"time"
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
	var pacer *Pacer
	if a.Pacer != nil {
		pacer = a.Pacer
	} else {
		pacer = globalPacer
	}

	start := time.Now()
	finite := a.Duration > 0

	for {
		select {
		case <-a.Context.Done():
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
