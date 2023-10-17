package metrics

import (
	"sync"
)

type DelegatingRegistry struct {
	*BaseRegistry
	registries     map[Registry]bool
	registriesLock sync.Mutex
}

func NewDelegatingRegistry() *DelegatingRegistry {
	registry := &DelegatingRegistry{
		registries: map[Registry]bool{},
	}
	registry.BaseRegistry = NewBaseRegistry(registry)
	return registry
}

func (r *DelegatingRegistry) NewCounter(id ID) Counter {
	counter := NewDelegatingCounter(id)
	r.addRegistries(counter)
	return counter
}

func (r *DelegatingRegistry) NewTimer(id ID) Timer {
	timer := NewDelegatingTimer(id)
	r.addRegistries(timer)
	return timer
}

func (r *DelegatingRegistry) Add(registry Registry) {
	if _, ok := registry.(*DelegatingRegistry); ok {
		return
	}

	r.registriesLock.Lock()
	defer r.registriesLock.Unlock()

	if _, exist := r.registries[registry]; exist {
		return
	}
	r.registries[registry] = true
	for _, metric := range r.Metrics() {
		if dm, ok := metric.(DelegatingMetric); ok {
			dm.Add(registry)
		}
	}
}

func (r *DelegatingRegistry) addRegistries(metric DelegatingMetric) {
	r.registriesLock.Lock()

	defer r.registriesLock.Unlock()
	for registry := range r.registries {
		metric.Add(registry)
	}
}

func (r *DelegatingRegistry) Close() {
	for registry := range r.registries {
		registry.Close()
	}
}
