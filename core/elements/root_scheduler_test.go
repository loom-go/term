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

		scheduler := NewScheduler(context.Background(), func() error {
			time.Sleep(10 * time.Millisecond)
			return nil
		}, make(chan error, 1))

		var updates []int

		wg.Go(func() {
			for i := range 100 {
				time.Sleep(2 * time.Millisecond)
				scheduler.Update(func() error {
					updates = append(updates, i)
					time.Sleep(time.Millisecond)
					return nil
				})
			}
		})

		wg.Wait()

		for i, update := range updates {
			assert.Equal(t, i, update)
		}
	})

	t.Run("recursive schedule", func(t *testing.T) {
		var wg sync.WaitGroup
		var logs []string

		wg.Add(1)
		scheduler := NewScheduler(context.Background(), func() error {
			time.Sleep(10 * time.Millisecond)
			logs = append(logs, "render")
			wg.Done()
			return nil
		}, make(chan error, 1))

		wg.Add(1)
		scheduler.Update(func() error {
			logs = append(logs, "first update")
			wg.Done()

			wg.Add(1)
			scheduler.Update(func() error {
				logs = append(logs, "second update")
				wg.Done()
				return nil
			})

			return nil
		})

		wg.Wait()

		assert.Equal(t, []string{"first update", "second update", "render"}, logs)
	})

	t.Run("access during flush", func(t *testing.T) {
		var logs []string

		scheduler := NewScheduler(context.Background(), func() error {
			time.Sleep(10 * time.Millisecond)
			logs = append(logs, "render")
			return nil
		}, make(chan error, 1))

		var wg sync.WaitGroup
		wg.Add(1)
		scheduler.Update(func() error {
			wg.Done()
			logs = append(logs, "update")
			return nil
		})
		wg.Wait()

		scheduler.Access(func() {
			logs = append(logs, "access")
		})

		// no need to wait for render. access should already block until flush is done

		assert.Equal(t, []string{"update", "render", "access"}, logs)
	})

	t.Run("update during heavy access", func(t *testing.T) {
		var wg sync.WaitGroup

		scheduler := NewScheduler(context.Background(), func() error {
			return nil
		}, make(chan error, 1))

		var updates []int

		wg.Go(func() {
			for i := range 100 {
				time.Sleep(2 * time.Millisecond)
				scheduler.Update(func() error {
					updates = append(updates, i)
					time.Sleep(time.Millisecond)
					return nil
				})
			}
		})

		nAccesses := 0
		wg.Go(func() {
			for len(updates) < 100 {
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
