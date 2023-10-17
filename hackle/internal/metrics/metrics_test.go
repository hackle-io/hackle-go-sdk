package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGlobalRegistry(t *testing.T) {
	assert.IsType(t, &DelegatingRegistry{}, globalRegistry)
}

func TestMetric(t *testing.T) {

	counter := NewCounter("counter", Tags{})
	timer := NewTimer("timer", Tags{})

	counter.Increment(1)
	timer.Record(time.Duration(42) * time.Millisecond)

	assert.Equal(t, int64(0), NewCounter("counter", Tags{}).Count())
	assert.Equal(t, int64(0), NewTimer("timer", Tags{}).Sum())

	cumulative := NewCumulativeRegistry()
	AddRegistry(cumulative)

	counter.Increment(1)
	timer.Record(time.Duration(1))

	assert.Equal(t, int64(1), NewCounter("counter", Tags{}).Count())
	assert.Equal(t, int64(1), NewTimer("timer", Tags{}).Sum())

	NewCounter("counter", Tags{"tag": "42"}).Increment(42)
	NewTimer("timer", Tags{"tag": "42"}).Record(time.Duration(42))

	assert.Equal(t, int64(1), NewCounter("counter", Tags{}).Count())
	assert.Equal(t, int64(42), NewCounter("counter", Tags{"tag": "42"}).Count())

	assert.Equal(t, int64(1), NewTimer("timer", Tags{}).Sum())
	assert.Equal(t, int64(42), NewTimer("timer", Tags{"tag": "42"}).Sum())
}
