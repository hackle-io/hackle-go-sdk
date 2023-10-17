package metrics

import (
	"sync/atomic"
	"time"
)

type CumulativeTimer struct {
	*BaseMetric
	*BaseTimer
	id    ID
	count int64
	sum   int64
	max   int64
}

func NewCumulativeTimer(id ID) *CumulativeTimer {
	timer := &CumulativeTimer{}
	timer.BaseMetric = NewBaseMetric(id)
	timer.BaseTimer = NewBaseTimer(timer)
	return timer
}

func (t *CumulativeTimer) Count() int64 {
	return atomic.LoadInt64(&t.count)
}

func (t *CumulativeTimer) Sum() int64 {
	return atomic.LoadInt64(&t.sum)
}

func (t *CumulativeTimer) Max() int64 {
	return atomic.LoadInt64(&t.max)
}

func (t *CumulativeTimer) Record(duration time.Duration) {
	if duration < 0 {
		return
	}
	atomic.AddInt64(&t.count, 1)
	atomic.AddInt64(&t.sum, int64(duration))
	t.updateMax(int64(duration))
}

func (t *CumulativeTimer) updateMax(value int64) {
	for {
		current := atomic.LoadInt64(&t.max)
		if current >= value {
			return
		}
		if atomic.CompareAndSwapInt64(&t.max, current, value) {
			return
		}
	}
}
