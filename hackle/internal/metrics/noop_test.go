package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNoopCounter(t *testing.T) {
	counter := NewNoopCounter(NewID("noop", Tags{}, TypeCounter))
	counter.Increment(42)
	assert.Equal(t, int64(0), counter.Count())
	assert.Equal(t, 0, len(counter.Measure()))
}

func TestNoopTimer(t *testing.T) {

	timer := NewNoopTimer(NewID("noop", Tags{}, TypeTimer))
	timer.Record(time.Duration(42))
	assert.Equal(t, int64(0), timer.Count())
	assert.Equal(t, int64(0), timer.Sum())
	assert.Equal(t, int64(0), timer.Max())
	assert.Equal(t, 0.0, timer.Mean())
	assert.Equal(t, 0, len(timer.Measure()))
}
