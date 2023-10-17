package metrics

import (
	"sync"
)

type DelegatingMetric interface {
	Metric
	Add(registry Registry)
	NoopMetric() Metric
	RegisterNewMetric(registry Registry) Metric
}

type BaseDelegatingMetric struct {
	DelegatingMetric
	registries     map[Registry]Metric
	registriesLock sync.Mutex
}

func NewBaseDelegatingMetric(self DelegatingMetric) *BaseDelegatingMetric {
	return &BaseDelegatingMetric{
		DelegatingMetric: self,
		registries:       make(map[Registry]Metric),
	}
}

func (m *BaseDelegatingMetric) Metrics() []Metric {
	m.registriesLock.Lock()
	defer m.registriesLock.Unlock()
	var ms []Metric
	for _, metric := range m.registries {
		ms = append(ms, metric)
	}
	return ms
}

func (m *BaseDelegatingMetric) First() Metric {
	ms := m.Metrics()
	if len(ms) == 0 {
		return m.NoopMetric()
	}
	return ms[0]
}

func (m *BaseDelegatingMetric) Add(registry Registry) {
	newMetric := m.RegisterNewMetric(registry)

	m.registriesLock.Lock()
	defer m.registriesLock.Unlock()

	newRegistries := make(map[Registry]Metric)
	for reg, metric := range m.registries {
		newRegistries[reg] = metric
	}
	newRegistries[registry] = newMetric
	m.registries = newRegistries
}
