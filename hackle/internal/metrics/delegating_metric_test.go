package metrics

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/concurrent"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelegatingMetric_Add(t *testing.T) {

	t.Run("concurrency", func(t *testing.T) {
		id := NewID("test", Tags{}, TypeCounter)
		delegatingMetric := NewDelegatingCounter(id)

		task := func() {
			for i := 0; i < 2_000; i++ {
				if i%2 == 0 {
					_ = delegatingMetric.Metrics()
				} else {
					delegatingMetric.Add(NewCumulativeRegistry())
				}
			}
		}

		executor := concurrent.NewExecutor()
		for i := 0; i < 4; i++ {
			executor.Go(task)
		}
		executor.Wait()

		assert.Equal(t, 4_000, len(delegatingMetric.Metrics()))
	})

	t.Run("add registry", func(t *testing.T) {
		id := NewID("test", Tags{}, TypeCounter)
		delegatingMetric := NewDelegatingCounter(id)

		for i := 0; i < 42; i++ {
			delegatingMetric.Add(NewCumulativeRegistry())
		}

		assert.Equal(t, 42, len(delegatingMetric.Metrics()))
	})
}

func TestDelegatingMetric_First(t *testing.T) {

	t.Run("when registry not added then return noop metric", func(t *testing.T) {
		sut := NewDelegatingCounter(NewID("counter", Tags{}, TypeCounter))
		metric := sut.First()
		assert.IsType(t, &NoopCounter{}, metric)
	})

	t.Run("when registry registered then return first registered metric", func(t *testing.T) {
		sut := NewDelegatingCounter(NewID("counter", Tags{}, TypeCounter))
		sut.Add(NewCumulativeRegistry())
		metric := sut.First()
		assert.IsType(t, &CumulativeCounter{}, metric)
	})
}
