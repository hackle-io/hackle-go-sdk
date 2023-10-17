package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDelegatingCounter_Count(t *testing.T) {
	t.Run("when registry not added then return 0", func(t *testing.T) {
		registry := NewDelegatingRegistry()
		counter := registry.Counter("counter", Tags{})
		assert.IsType(t, &DelegatingCounter{}, counter)
		assert.IsType(t, int64(0), counter.Count())
	})

}

func TestDelegatingCounter_Increment(t *testing.T) {
	t.Run("when registry added then delegate to registered metric", func(t *testing.T) {
		delegating := NewDelegatingRegistry()
		cumulative1 := NewCumulativeRegistry()
		cumulative2 := NewCumulativeRegistry()

		delegating.Add(cumulative1)
		delegating.Add(cumulative2)

		counter := delegating.Counter("counter", Tags{})
		counter.Increment(42)

		assert.Equal(t, int64(42), counter.Count())
		assert.Equal(t, int64(42), delegating.Counter("counter", Tags{}).Count())
		assert.Equal(t, int64(42), cumulative1.Counter("counter", Tags{}).Count())
		assert.Equal(t, int64(42), cumulative2.Counter("counter", Tags{}).Count())
	})
}

func TestDelegatingCounter_Measure(t *testing.T) {
	delegating := NewDelegatingRegistry()
	counter := delegating.Counter("counter", Tags{})
	measurements := counter.Measure()

	assert.Equal(t, 1, len(measurements))
	assert.Equal(t, 0.0, measurements[0].Value())

	counter.Increment(42)
	assert.Equal(t, 0.0, measurements[0].Value())

	delegating.Add(NewCumulativeRegistry())
	counter.Increment(42)
	assert.Equal(t, 42.0, measurements[0].Value())
}
