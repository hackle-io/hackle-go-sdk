package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/concurrent"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter_Increment(t *testing.T) {

	t.Run("increment only", func(t *testing.T) {
		c := counter()
		for i := 0; i < 100; i++ {
			c.Increment(1)
			assert.Equal(t, int64(i+1), c.Count())
		}
	})

	t.Run("with flush", func(t *testing.T) {
		c := counter()

		c.Increment(1)
		assert.Equal(t, int64(1), c.Count())

		fc := c.Flush()
		assert.Equal(t, int64(0), c.Count())
		assert.Equal(t, int64(1), fc.(metrics.Counter).Count())

		c.Increment(42)
		assert.Equal(t, int64(42), c.Count())
	})

	t.Run("concurrency", func(t *testing.T) {
		c := counter()

		counters := make([][]metrics.Metric, 16)

		task := func(n int) {
			flushed := make([]metrics.Metric, 0)
			for i := 0; i < 100_000; i++ {
				if i%2 == 0 {
					c.Increment(2)
				} else {
					flushed = append(flushed, c.Flush())
				}
			}
			counters[n] = flushed
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
		for _, cs := range counters {
			for _, metric := range cs {
				count += metric.(metrics.Counter).Count()
			}
		}
		count += c.Flush().(metrics.Counter).Count()
		assert.Equal(t, int64(1_600_000), count)
	})
}

func TestCounter_Measure(t *testing.T) {
	c := counter()
	c.Increment(42)
	measurements := c.Measure()

	assert.Equal(t, 42.0, measurements[0].Value())

	c.Flush()
	assert.Equal(t, 0.0, measurements[0].Value())
}

func counter() *Counter {
	return NewCounter(metrics.NewID("counter", metrics.Tags{}, metrics.TypeCounter))
}
