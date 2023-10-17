package metrics

import (
	"sync/atomic"
)

type CumulativeCounter struct {
	*BaseMetric
	*BaseCounter
	value int64
}

func NewCumulativeCounter(id ID) *CumulativeCounter {
	counter := &CumulativeCounter{}
	counter.BaseMetric = NewBaseMetric(id)
	counter.BaseCounter = NewBaseCounter(counter)
	return counter
}

func (c *CumulativeCounter) Count() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *CumulativeCounter) Increment(delta int64) {
	atomic.AddInt64(&c.value, delta)
}
