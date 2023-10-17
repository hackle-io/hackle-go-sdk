package metrics

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/concurrent"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCounter_Increment(t *testing.T) {

	t.Run("increment", func(t *testing.T) {
		counter := NewCumulativeRegistry().Counter("counter", Tags{})

		task := func() {
			for i := 0; i < 100_000; i++ {
				counter.Increment(1)
			}
		}

		executor := concurrent.NewExecutor()
		for i := 0; i < 8; i++ {
			executor.Go(task)
		}
		executor.Wait()

		assert.Equal(t, int64(800_000), counter.Count())
	})
}

func TestCounter_Measure(t *testing.T) {
	registry := NewCumulativeRegistry()
	counter := registry.Counter("counter", Tags{})
	measurements := counter.Measure()

	assert.Equal(t, 1, len(measurements))
	assert.Equal(t, 0.0, measurements[0].Value())

	counter.Increment(42)
	assert.Equal(t, 42.0, measurements[0].Value())
}
