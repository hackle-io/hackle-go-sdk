package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDelegatingTimer_Get(t *testing.T) {
	t.Run("when registry not added then return 0 ", func(t *testing.T) {
		timer := NewDelegatingRegistry().Timer("timer", Tags{})
		assert.IsType(t, &DelegatingTimer{}, timer)

		assert.Equal(t, int64(0), timer.Count())
		assert.Equal(t, int64(0), timer.Sum())
		assert.Equal(t, int64(0), timer.Max())
		assert.Equal(t, 0.0, timer.Mean())
	})
}

func TestDelegatingTimer_Record(t *testing.T) {
	t.Run("delegate to registered metric", func(t *testing.T) {
		delegating := NewDelegatingRegistry()
		cumulative1 := NewCumulativeRegistry()
		cumulative2 := NewCumulativeRegistry()

		delegating.Add(cumulative1)
		delegating.Add(cumulative2)

		timer := delegating.Timer("timer", Tags{})
		timer.Record(time.Duration(42))

		assert.Equal(t, int64(42), timer.Sum())
		assert.Equal(t, int64(42), delegating.Timer("timer", Tags{}).Sum())
		assert.Equal(t, int64(42), cumulative1.Timer("timer", Tags{}).Sum())
		assert.Equal(t, int64(42), cumulative2.Timer("timer", Tags{}).Sum())
	})
}

func TestDelegatingTimer_Measure(t *testing.T) {
	registry := NewDelegatingRegistry()
	timer := registry.Timer("timer", Tags{})
	measurements := timer.Measure()

	assert.Equal(t, 4, len(measurements))
	assert.Equal(t, 0.0, measurements[0].Value())
	assert.Equal(t, 0.0, measurements[1].Value())
	assert.Equal(t, 0.0, measurements[2].Value())
	assert.Equal(t, 0.0, measurements[3].Value())

	timer.Record(time.Duration(42 * time.Millisecond))
	assert.Equal(t, 0.0, measurements[0].Value())
	assert.Equal(t, 0.0, measurements[1].Value())
	assert.Equal(t, 0.0, measurements[2].Value())
	assert.Equal(t, 0.0, measurements[3].Value())

	registry.Add(NewCumulativeRegistry())
	timer.Record(time.Duration(42 * time.Millisecond))
	assert.Equal(t, 1.0, measurements[0].Value())
	assert.Equal(t, 42.0, measurements[1].Value())
	assert.Equal(t, 42.0, measurements[2].Value())
	assert.Equal(t, 42.0, measurements[3].Value())
}
