package metrics

import "time"

type NoopCounter struct {
	*BaseMetric
}

func NewNoopCounter(id ID) *NoopCounter {
	return &NoopCounter{
		BaseMetric: NewBaseMetric(id),
	}
}

func (c *NoopCounter) Measure() []Measurement {
	return []Measurement{}
}

func (c *NoopCounter) Count() int64 {
	return 0
}

func (c *NoopCounter) Increment(int64) {}

type NoopTimer struct {
	*BaseMetric
}

func NewNoopTimer(id ID) *NoopTimer {
	return &NoopTimer{
		BaseMetric: NewBaseMetric(id),
	}
}

func (t *NoopTimer) Measure() []Measurement {
	return []Measurement{}
}

func (t *NoopTimer) Count() int64 {
	return 0
}

func (t *NoopTimer) Sum() int64 {
	return 0
}

func (t *NoopTimer) Max() int64 {
	return 0
}

func (t *NoopTimer) Mean() float64 {
	return 0
}

func (t *NoopTimer) Record(time.Duration) {}
