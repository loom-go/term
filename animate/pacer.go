package animate

import (
	"sync"
	"time"

	"github.com/loom-go/loom/signals"
)

const defaultFPS = 60

var globalPacer = NewPacer(time.Second / defaultFPS)

func Pace(tick func(time.Time)) {
	globalPacer.Pace(tick)
}

type frameRequest struct {
	tick func(now time.Time)
	done chan struct{}
}

type Pacer struct {
	mu       sync.Mutex
	rate     time.Duration
	requests []*frameRequest
}

// NewPacer creates a new Pacer that paces frame updates at the given rate (frames per second).
// It can be given to animate.A to control the pacing of animations.
//
// By default, animations use a global FramePacer at 60 FPS.
func NewPacer(rate time.Duration) *Pacer {
	p := &Pacer{
		rate:     rate,
		requests: make([]*frameRequest, 0),
	}

	go p.loop()
	return p
}

func (p *Pacer) loop() {
	ticker := time.NewTicker(p.rate)
	defer ticker.Stop()

	for now := range ticker.C {
		p.mu.Lock()
		reqs := p.requests
		p.requests = nil
		p.mu.Unlock()

		if len(reqs) == 0 {
			continue
		}

		signals.Batch(func() {
			for _, req := range reqs {
				req.tick(now)
				close(req.done)
			}
		})
	}
}

func (p *Pacer) Pace(tick func(time.Time)) {
	req := &frameRequest{
		tick: tick,
		done: make(chan struct{}),
	}

	p.mu.Lock()
	p.requests = append(p.requests, req)
	p.mu.Unlock()

	<-req.done
}
