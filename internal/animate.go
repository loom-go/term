package internal

import (
	"time"
)

// todo: easing
type Animation struct {
	Duration time.Duration
	Tick     func(progress float64)
}

func Animate(a Animation) {
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
