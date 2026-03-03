package elements

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Run("concurrent fifo", func(t *testing.T) {
		var wg sync.WaitGroup

		scheduler := NewScheduler(context.Background())

		var updates []int

		wg.Go(func() {
			for i := range 100 {
				time.Sleep(2 * time.Millisecond)
				scheduler.Schedule(taskUpdate, func() error {
					updates = append(updates, i)
					time.Sleep(time.Millisecond)
					return nil
				})
			}
		})

		wg.Go(func() {
			for len(updates) < 100 {
				<-scheduler.Schedule(taskRender, func() error {
					time.Sleep(5 * time.Millisecond)
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

		scheduler := NewScheduler(context.Background())

		var logs []string

		wg.Add(1)
		scheduler.Schedule(taskUpdate, func() error {
			logs = append(logs, "first update")
			wg.Done()

			wg.Add(1)
			scheduler.Schedule(taskRender, func() error {
				logs = append(logs, "render")
				wg.Done()
				return nil
			})

			wg.Add(1)
			scheduler.Schedule(taskUpdate, func() error {
				logs = append(logs, "second update")
				wg.Done()
				return nil
			})

			return nil
		})

		wg.Wait()

		assert.Equal(t, []string{"first update", "render", "second update"}, logs)
	})
}
