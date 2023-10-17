package metrics

type DelegatingCounter struct {
	*BaseMetric
	*BaseCounter
	*BaseDelegatingMetric
	id ID
}

func NewDelegatingCounter(id ID) *DelegatingCounter {
	counter := &DelegatingCounter{}
	counter.BaseMetric = NewBaseMetric(id)
	counter.BaseCounter = NewBaseCounter(counter)
	counter.BaseDelegatingMetric = NewBaseDelegatingMetric(counter)
	return counter
}

func (c *DelegatingCounter) Count() int64 {
	return c.First().(Counter).Count()
}

func (c *DelegatingCounter) Increment(delta int64) {
	for _, metric := range c.Metrics() {
		metric.(Counter).Increment(delta)
	}
}

func (c *DelegatingCounter) NoopMetric() Metric {
	return NewNoopCounter(c.ID())
}

func (c *DelegatingCounter) RegisterNewMetric(registry Registry) Metric {
	return registry.RegisterCounter(c.ID())
}
