package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"sync"
)

type Metric interface {
	metrics.Metric
	InitialMetric() metrics.Metric
	Current() metrics.Metric
	Flush() metrics.Metric
}

type BaseFlushMetric struct {
	Metric
	current metrics.Metric
	mu      sync.RWMutex
}

func NewBaseMetric(self Metric) *BaseFlushMetric {
	return &BaseFlushMetric{
		Metric:  self,
		current: self.InitialMetric(),
	}
}

func (m *BaseFlushMetric) Current() metrics.Metric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.current
}

func (m *BaseFlushMetric) Flush() metrics.Metric {
	m.mu.Lock()
	defer m.mu.Unlock()
	current := m.current
	m.current = m.InitialMetric()
	return current
}
