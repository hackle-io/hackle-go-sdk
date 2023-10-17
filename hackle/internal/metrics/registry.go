package metrics

import (
	"fmt"
	"sync"
)

type Registry interface {
	Metrics() []Metric

	Counter(name string, tags Tags) Counter
	Timer(name string, tags Tags) Timer

	RegisterCounter(id ID) Counter
	RegisterTimer(id ID) Timer

	NewCounter(id ID) Counter
	NewTimer(id ID) Timer

	Close()
}

type BaseRegistry struct {
	Registry
	metrics map[string]Metric
	lock    sync.Mutex
}

func NewBaseRegistry(self Registry) *BaseRegistry {
	return &BaseRegistry{
		Registry: self,
		metrics:  map[string]Metric{},
	}
}

func (r *BaseRegistry) Metrics() []Metric {
	metrics := make([]Metric, 0, len(r.metrics))
	for _, metric := range r.metrics {
		metrics = append(metrics, metric)
	}
	return metrics
}

func (r *BaseRegistry) Counter(name string, tags Tags) Counter {
	return NewCounterBuilder(name).Tags(tags).Register(r)
}

func (r *BaseRegistry) Timer(name string, tags Tags) Timer {
	return NewTimerBuilder(name).Tags(tags).Register(r)
}

func (r *BaseRegistry) RegisterCounter(id ID) Counter {
	metric := r.registerMetricIfNecessary(id, func(id ID) Metric {
		return r.NewCounter(id)
	})
	if counter, ok := metric.(Counter); ok {
		return counter
	}
	panic(fmt.Sprintf("metric already registered with different type: %T, %s (expected: metrics.Counter)", id.fqn, metric))
}

func (r *BaseRegistry) RegisterTimer(id ID) Timer {
	metric := r.registerMetricIfNecessary(id, func(id ID) Metric {
		return r.NewTimer(id)
	})
	if timer, ok := metric.(Timer); ok {
		return timer
	}
	panic(fmt.Sprintf("metric already registered with different type: %T, %s (expected: metrics.Counter)", id.fqn, metric))
}

func (r *BaseRegistry) registerMetricIfNecessary(id ID, create func(ID) Metric) Metric {
	r.lock.Lock()
	defer r.lock.Unlock()
	if metric, ok := r.metrics[id.fqn]; ok {
		return metric
	}
	metric := create(id)
	r.metrics[id.fqn] = metric
	return metric
}
