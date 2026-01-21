package debug

import (
	"context"
	"slices"
	"sync"
	"time"
)

type TimingRecord struct {
	Last time.Duration
	Avg  time.Duration
	Min  time.Duration
	Max  time.Duration
}

type TimingMetric struct {
	mu sync.Mutex

	emitter *Emitter[*TimingRecord]

	count int64
	min   time.Duration
	max   time.Duration
	total time.Duration
}

func NewTimingMetric() *TimingMetric {
	return &TimingMetric{emitter: NewEmitter[*TimingRecord]()}
}

func (s *TimingMetric) Emit(d time.Duration) {
	s.mu.Lock()
	s.count++
	s.total += d

	if s.min == 0 || d < s.min {
		s.min = d
	}
	if d > s.max {
		s.max = d
	}

	record := &TimingRecord{
		Last: d,
		Avg:  s.total / time.Duration(s.count),
		Min:  s.min,
		Max:  s.max,
	}
	s.mu.Unlock()

	s.emitter.Emit(record)
}

func (s *TimingMetric) Subscribe(buffer int) (ch <-chan *TimingRecord, cancel func()) {
	ctx, cancel := context.WithCancel(context.Background())
	return s.emitter.Subscribe(ctx, buffer), cancel
}

func (s *TimingMetric) Reset() {
	s.mu.Lock()
	s.count = 0
	s.min = 0
	s.max = 0
	s.total = 0
	s.mu.Unlock()
}

type RateMetric struct {
	mu sync.Mutex

	emitter *Emitter[float64]

	samples []time.Time
	window  time.Duration
}

func NewRateMetric(window time.Duration) *RateMetric {
	return &RateMetric{
		emitter: NewEmitter[float64](),
		window:  window,
	}
}

func (s *RateMetric) Emit() {
	s.mu.Lock()
	now := time.Now()
	cutoff := now.Add(-s.window)

	s.samples = append(s.samples, now)
	s.samples = slices.DeleteFunc(s.samples, func(t time.Time) bool {
		return t.Before(cutoff)
	})

	var rate float64
	if len(s.samples) >= 2 {
		// Use actual span from oldest to newest sample
		actualSpan := s.samples[len(s.samples)-1].Sub(s.samples[0]).Seconds()
		if actualSpan > 0 {
			// N samples have N-1 intervals between them
			rate = float64(len(s.samples)-1) / actualSpan
		}
	}
	s.mu.Unlock()

	s.emitter.Emit(rate)
}

func (s *RateMetric) Subscribe(buffer int) (ch <-chan float64, cancel func()) {
	ctx, cancel := context.WithCancel(context.Background())
	return s.emitter.Subscribe(ctx, buffer), cancel
}
