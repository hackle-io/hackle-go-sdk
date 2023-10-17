package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/concurrent"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimer_Record(t *testing.T) {

	t.Run("record only", func(t *testing.T) {

		timer := newTimer()
		for i := 0; i < 100; i++ {
			timer.Record(time.Duration(i + 1))
		}
		assert.Equal(t, int64(100), timer.Count())
		assert.Equal(t, int64(5050), timer.Sum())
		assert.Equal(t, int64(100), timer.Max())
		assert.Equal(t, 50.5, timer.Mean())
	})

	t.Run("with flush", func(t *testing.T) {
		timer := newTimer()
		for i := 0; i < 100; i++ {
			timer.Record(time.Duration(i + 1))
		}
		assert.Equal(t, int64(100), timer.Count())
		assert.Equal(t, int64(5050), timer.Sum())
		assert.Equal(t, int64(100), timer.Max())
		assert.Equal(t, 50.5, timer.Mean())

		flush := timer.Flush()
		assert.Equal(t, int64(0), timer.Count())
		assert.Equal(t, int64(0), timer.Sum())
		assert.Equal(t, int64(0), timer.Max())
		assert.Equal(t, 0.0, timer.Mean())

		flushedTimer := flush.(metrics.Timer)
		assert.Equal(t, int64(100), flushedTimer.Count())
		assert.Equal(t, int64(5050), flushedTimer.Sum())
		assert.Equal(t, int64(100), flushedTimer.Max())
		assert.Equal(t, 50.5, flushedTimer.Mean())
	})
	t.Run("concurrency", func(t *testing.T) {
		timer := newTimer()

		timers := make([][]metrics.Metric, 16)

		task := func(n int) {
			flushed := make([]metrics.Metric, 0)
			for i := 0; i < 100_000; i++ {
				if i%2 == 0 {
					timer.Record(time.Duration(2))
				} else {
					flushed = append(flushed, timer.Flush())
				}
			}
			timers[n] = flushed
		}

		executor := concurrent.NewExecutor()
		for i := 0; i < 16; i++ {
			n := i
			executor.Go(func() {
				task(n)
			})
		}
		executor.Wait()

		var count int64
		var sum int64

		for _, ts := range timers {
			for _, metric := range ts {
				t := metric.(metrics.Timer)
				count += t.Count()
				sum += t.Sum()
			}
		}

		assert.Equal(t, int64(800_000), count)
		assert.Equal(t, int64(1_600_000), sum)
	})
}

func newTimer() *Timer {
	return NewTimer(metrics.NewID("timer", metrics.Tags{}, metrics.TypeTimer))
}
