package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics/push"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"time"
)

type MetricRegistry interface {
	Flush(ms []metrics.Metric)
}

type BaseMetricRegistry struct {
	MetricRegistry
	*metrics.BaseRegistry
	*push.BaseMetricRegistry
}

func NewBaseMetricRegistry(self MetricRegistry, scheduler schedule.Scheduler, pushInterval time.Duration) *BaseMetricRegistry {
	registry := &BaseMetricRegistry{
		MetricRegistry: self,
	}
	registry.BaseRegistry = metrics.NewBaseRegistry(registry)
	registry.BaseMetricRegistry = push.NewBaseMetricRegistry(registry, scheduler, pushInterval)
	return registry
}

func (r *BaseMetricRegistry) NewCounter(id metrics.ID) metrics.Counter {
	return NewCounter(id)
}

func (r *BaseMetricRegistry) NewTimer(id metrics.ID) metrics.Timer {
	return NewTimer(id)
}

func (r *BaseMetricRegistry) Close() {
	r.Stop()
}

func (r *BaseMetricRegistry) Publish() {
	ms := make([]metrics.Metric, 0)
	for _, metric := range r.Metrics() {
		if flushMetric, ok := metric.(Metric); ok {
			ms = append(ms, flushMetric.Flush())

		}
	}
	r.Flush(ms)
}
