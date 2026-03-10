package elements

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Run("fifo updates", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var updates []int

		opts := SchedulerOptions{
			Context: context.Background(),
			Render: func() error {
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		}
		scheduler := NewScheduler(opts)

		wg.Go(func() {
			for i := range 100 {
				time.Sleep(2 * time.Millisecond)
				scheduler.Update(func() error {
					mu.Lock()
					updates = append(updates, i)
					mu.Unlock()
					time.Sleep(time.Millisecond)
					return nil
				})
			}
		})

		wg.Wait()

		mu.Lock()
		for i, update := range updates {
			assert.Equal(t, i, update)
		}
		mu.Unlock()
	})

	t.Run("recursive schedule", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var logs []string

		opts := SchedulerOptions{
			Context: context.Background(),
			Render: func() error {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				logs = append(logs, "render")
				mu.Unlock()
				wg.Done()
				return nil
			},
		}

		wg.Add(1)
		scheduler := NewScheduler(opts)

		wg.Add(1)
		scheduler.Update(func() error {
			mu.Lock()
			logs = append(logs, "first update")
			mu.Unlock()
			wg.Done()

			wg.Add(1)
			scheduler.Update(func() error {
				mu.Lock()
				logs = append(logs, "second update")
				mu.Unlock()
				wg.Done()
				return nil
			})

			return nil
		})

		wg.Wait()

		assert.Equal(t, []string{"first update", "second update", "render"}, logs)
	})

	t.Run("access during flush", func(t *testing.T) {
		var mu sync.Mutex
		var logs []string

		opts := SchedulerOptions{
			Context: context.Background(),
			Render: func() error {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				logs = append(logs, "render")
				mu.Unlock()
				return nil
			},
		}
		scheduler := NewScheduler(opts)

		var wg sync.WaitGroup
		wg.Add(1)
		scheduler.Update(func() error {
			wg.Done()
			mu.Lock()
			logs = append(logs, "update")
			mu.Unlock()
			return nil
		})
		wg.Wait()

		scheduler.Access(func() {
			mu.Lock()
			logs = append(logs, "access")
			mu.Unlock()
		})

		// no need to wait for render. access should already block until flush is done

		mu.Lock()
		assert.Equal(t, []string{"update", "render", "access"}, logs)
		mu.Unlock()
	})

	t.Run("update during heavy access", func(t *testing.T) {
		var wg sync.WaitGroup
		var mu sync.Mutex
		var updates []int

		opts := SchedulerOptions{
			Context: context.Background(),
			Render: func() error {
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		}
		scheduler := NewScheduler(opts)

		wg.Go(func() {
			for i := range 100 {
				time.Sleep(2 * time.Millisecond)
				scheduler.Update(func() error {
					mu.Lock()
					updates = append(updates, i)
					mu.Unlock()
					time.Sleep(time.Millisecond)
					return nil
				})
			}
		})

		nAccesses := 0
		wg.Go(func() {

			for {
				mu.Lock()
				if len(updates) < 100 {
					mu.Unlock()
					break
				}
				mu.Unlock()

				time.Sleep(time.Millisecond)
				scheduler.Access(func() {
					time.Sleep(time.Millisecond)
					nAccesses++
				})
			}
		})

		wg.Wait()

		// make sure accesses can't stop updates from being flushed for too long
		assert.Less(t, nAccesses, 110)
	})
}
