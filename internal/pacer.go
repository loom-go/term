package internal

import (
	"sync"
	"time"

	"github.com/AnatoleLucet/loom/signals"
)

var pacer = NewFramePacer()

const frameRate = 60
const frameDuration = time.Second / frameRate

type frameRequest struct {
	tick func(now time.Time)
	done chan struct{}
}

type FramePacer struct {
	mu       sync.Mutex
	requests []*frameRequest
}

func NewFramePacer() *FramePacer {
	p := &FramePacer{
		requests: make([]*frameRequest, 0),
	}

	go p.loop()
	return p
}

func (p *FramePacer) loop() {
	ticker := time.NewTicker(frameDuration)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		reqs := p.requests
		p.requests = p.requests[:0]
		p.mu.Unlock()

		if len(reqs) == 0 {
			continue
		}

		signals.Batch(func() {
			now := time.Now()
			for _, req := range reqs {
				defer close(req.done)
				req.tick(now)
			}
		})
	}
}

func (p *FramePacer) Pace(tick func(time.Time)) {
	req := &frameRequest{
		tick: tick,
		done: make(chan struct{}),
	}

	p.mu.Lock()
	p.requests = append(p.requests, req)
	p.mu.Unlock()

	<-req.done
}
