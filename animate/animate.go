package animate

import (
	"time"
)

// todo: easing and infinite p0..p1..p0..p1.. (currently infinite is just locked at p0)

// A represents an animation that can be run with Run.
type A struct {
	Duration time.Duration
	Tick     func(progress float64)
	Pacer    *FramePacer
}

// Run executes the given animation A and blocks until it is complete.
func Run(a A) {
	var pacer *FramePacer
	if a.Pacer != nil {
		pacer = a.Pacer
	} else {
		pacer = globalPacer
	}

	start := time.Now()
	finite := a.Duration > 0

	for {
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
