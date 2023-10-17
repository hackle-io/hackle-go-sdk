package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"time"
)

type Timer struct {
	*metrics.BaseMetric
	*metrics.BaseTimer
	*BaseFlushMetric
}

func NewTimer(id metrics.ID) *Timer {
	timer := &Timer{}
	timer.BaseMetric = metrics.NewBaseMetric(id)
	timer.BaseTimer = metrics.NewBaseTimer(timer)
	timer.BaseFlushMetric = NewBaseMetric(timer)
	return timer
}

func (t *Timer) Count() int64 {
	return t.Current().(metrics.Timer).Count()
}

func (t *Timer) Sum() int64 {
	return t.Current().(metrics.Timer).Sum()
}

func (t *Timer) Max() int64 {
	return t.Current().(metrics.Timer).Max()
}

func (t *Timer) Record(duration time.Duration) {
	t.Current().(metrics.Timer).Record(duration)
}

func (t *Timer) InitialMetric() metrics.Metric {
	return metrics.NewCumulativeTimer(t.ID())
}
