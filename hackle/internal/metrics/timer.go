package metrics

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/times"
	"time"
)

type Timer interface {
	Metric
	Count() int64
	Sum() int64
	Max() int64
	Mean() float64
	Record(duration time.Duration)
}

type BaseTimer struct {
	Timer
}

func NewBaseTimer(self Timer) *BaseTimer {
	return &BaseTimer{
		Timer: self,
	}
}

func (t *BaseTimer) Mean() float64 {
	count := t.Count()
	if count == 0 {
		return 0
	}
	return float64(t.Sum()) / float64(count)
}

func (t *BaseTimer) Measure() []Measurement {
	return []Measurement{
		NewMeasurement(FieldCount, func() float64 { return float64(t.Count()) }),
		NewMeasurement(FieldTotal, func() float64 { return times.Millis(float64(t.Sum())) }),
		NewMeasurement(FieldMax, func() float64 { return times.Millis(float64(t.Max())) }),
		NewMeasurement(FieldMean, func() float64 { return times.Millis(t.Mean()) }),
	}
}

type TimerBuilder struct {
	name string
	tags Tags
}

func NewTimerBuilder(name string) *TimerBuilder {
	return &TimerBuilder{
		name: name,
		tags: map[string]string{},
	}
}

func (b *TimerBuilder) Tag(key string, value string) *TimerBuilder {
	b.tags[key] = value
	return b
}

func (b *TimerBuilder) Tags(tags Tags) *TimerBuilder {
	for key, value := range tags {
		b.Tag(key, value)
	}
	return b
}

func (b *TimerBuilder) Register(registry Registry) Timer {
	id := NewID(b.name, b.tags, TypeTimer)
	return registry.RegisterTimer(id)
}

func NewTimerSample(c clock.Clock) *TimerSample {
	return &TimerSample{
		clock:     c,
		startTick: c.Tick(),
	}
}

type TimerSample struct {
	clock     clock.Clock
	startTick int64
}

func (s *TimerSample) Stop(timer Timer) time.Duration {
	duration := time.Duration(s.clock.Tick() - s.startTick)
	timer.Record(duration)
	return duration
}
