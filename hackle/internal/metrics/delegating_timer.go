package metrics

import "time"

type DelegatingTimer struct {
	*BaseMetric
	*BaseTimer
	*BaseDelegatingMetric
}

func NewDelegatingTimer(id ID) *DelegatingTimer {
	timer := &DelegatingTimer{}
	timer.BaseMetric = NewBaseMetric(id)
	timer.BaseTimer = NewBaseTimer(timer)
	timer.BaseDelegatingMetric = NewBaseDelegatingMetric(timer)
	return timer
}

func (t *DelegatingTimer) Count() int64 {
	return t.First().(Timer).Count()
}

func (t *DelegatingTimer) Sum() int64 {
	return t.First().(Timer).Sum()
}

func (t *DelegatingTimer) Max() int64 {
	return t.First().(Timer).Max()
}

func (t *DelegatingTimer) Record(duration time.Duration) {
	for _, metric := range t.Metrics() {
		metric.(Timer).Record(duration)
	}
}

func (t *DelegatingTimer) NoopMetric() Metric {
	return NewNoopTimer(t.ID())
}

func (t *DelegatingTimer) RegisterNewMetric(registry Registry) Metric {
	return registry.RegisterTimer(t.ID())
}
