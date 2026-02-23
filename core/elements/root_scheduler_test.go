package elements

import (
	"context"
	"fmt"
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
					fmt.Printf("Scheduled update %d\n", i)
					return nil
				})
			}
		})

		wg.Go(func() {
			for len(updates) < 100 {
				<-scheduler.Schedule(taskRender, func() error {
					time.Sleep(5 * time.Millisecond)
					fmt.Println("Scheduled render")
					return nil
				})
			}
		})

		wg.Wait()

		for i, update := range updates {
			assert.Equal(t, i, update)
		}
	})
}
