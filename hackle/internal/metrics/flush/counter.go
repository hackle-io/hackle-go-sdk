package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
)

type Counter struct {
	*metrics.BaseMetric
	*metrics.BaseCounter
	*BaseFlushMetric
}

func NewCounter(id metrics.ID) *Counter {
	counter := &Counter{}
	counter.BaseMetric = metrics.NewBaseMetric(id)
	counter.BaseCounter = metrics.NewBaseCounter(counter)
	counter.BaseFlushMetric = NewBaseMetric(counter)
	return counter
}

func (c *Counter) Count() int64 {
	return c.Current().(metrics.Counter).Count()
}

func (c *Counter) Increment(delta int64) {
	c.Current().(metrics.Counter).Increment(delta)
}

func (c *Counter) InitialMetric() metrics.Metric {
	return metrics.NewCumulativeCounter(c.ID())
}
