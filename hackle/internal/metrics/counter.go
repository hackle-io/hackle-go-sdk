package metrics

type Counter interface {
	Metric
	Count() int64
	Increment(delta int64)
}

type BaseCounter struct {
	Counter
}

func NewBaseCounter(self Counter) *BaseCounter {
	return &BaseCounter{
		Counter: self,
	}
}

func (c *BaseCounter) Measure() []Measurement {
	return []Measurement{
		NewMeasurement(FieldCount, func() float64 { return float64(c.Count()) }),
	}
}

type CounterBuilder struct {
	name string
	tags Tags
}

func NewCounterBuilder(name string) *CounterBuilder {
	return &CounterBuilder{
		name: name,
		tags: map[string]string{},
	}
}

func (b *CounterBuilder) Tag(key string, value string) *CounterBuilder {
	b.tags[key] = value
	return b
}

func (b *CounterBuilder) Tags(tags Tags) *CounterBuilder {
	for key, value := range tags {
		b.Tag(key, value)
	}
	return b
}

func (b *CounterBuilder) Register(registry Registry) Counter {
	id := NewID(b.name, b.tags, TypeCounter)
	return registry.RegisterCounter(id)
}
