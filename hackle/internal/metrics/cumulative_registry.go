package metrics

type CumulativeRegistry struct {
	*BaseRegistry
}

func NewCumulativeRegistry() *CumulativeRegistry {
	registry := &CumulativeRegistry{}
	registry.BaseRegistry = NewBaseRegistry(registry)
	return registry
}

func (r *CumulativeRegistry) NewCounter(id ID) Counter {
	return NewCumulativeCounter(id)
}

func (r *CumulativeRegistry) NewTimer(id ID) Timer {
	return NewCumulativeTimer(id)
}

func (r *CumulativeRegistry) Close() {}
