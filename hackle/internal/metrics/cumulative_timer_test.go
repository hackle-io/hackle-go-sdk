package metrics

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/concurrent"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimer_Record(t *testing.T) {

	t.Run("negative duration should be ignored", func(t *testing.T) {
		timer := NewCumulativeRegistry().Timer("timer", Tags{})
		timer.Record(time.Duration(-1))
		assert.Equal(t, int64(0), timer.Count())
	})

	t.Run("concurrency", func(t *testing.T) {

		timer := NewCumulativeRegistry().Timer("timer", Tags{})

		task := func() {
			for i := 0; i < 100_000; i++ {
				timer.Record(time.Duration(i + 1))
			}
		}

		executor := concurrent.NewExecutor()
		for i := 0; i < 8; i++ {
			executor.Go(task)
		}
		executor.Wait()

		assert.Equal(t, int64(800_000), timer.Count())
		assert.Equal(t, int64(40000400000), timer.Sum())
		assert.Equal(t, int64(100_000), timer.Max())
		assert.Equal(t, 50_000.5, timer.Mean())
	})
}
