package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBaseMetricRegistry_NewCounter(t *testing.T) {
	registry := newMockFlushMetricRegistry()
	counter := registry.NewCounter(metrics.NewID("counter", metrics.Tags{}, metrics.TypeCounter))
	assert.IsType(t, &Counter{}, counter)
}
func TestBaseMetricRegistry_NewTimer(t *testing.T) {
	registry := newMockFlushMetricRegistry()
	timer := registry.NewTimer(metrics.NewID("timer", metrics.Tags{}, metrics.TypeTimer))
	assert.IsType(t, &Timer{}, timer)
}

func TestBaseMetricRegistry_Publish(t *testing.T) {
	registry := newMockFlushMetricRegistry()

	registry.Counter("counter", metrics.Tags{})
	registry.Counter("counter", metrics.Tags{"tag-1": "tag-2"})
	registry.Timer("timer", metrics.Tags{})

	registry.Publish()
	assert.Equal(t, 3, len(registry.flushed))
	assert.Equal(t, 3, len(registry.Metrics()))

	registry.Timer("timer", metrics.Tags{"tag-1": "tag-2"})
	registry.Publish()
	assert.Equal(t, 7, len(registry.flushed))
}
func TestBaseMetricRegistry_Close(t *testing.T) {
	registry := newMockFlushMetricRegistry()
	registry.Counter("counter", metrics.Tags{})
	registry.Start()

	assert.Equal(t, 0, len(registry.flushed))

	registry.Close()
	assert.Equal(t, 1, len(registry.flushed))
}

func newMockFlushMetricRegistry() *mockFlushMetricRegistry {
	registry := &mockFlushMetricRegistry{}
	registry.BaseMetricRegistry = NewBaseMetricRegistry(registry, schedule.NewTickerScheduler(), 1*time.Hour)
	return registry
}

type mockFlushMetricRegistry struct {
	*BaseMetricRegistry
	flushed []metrics.Metric
}

func (m *mockFlushMetricRegistry) Flush(ms []metrics.Metric) {
	m.flushed = append(m.flushed, ms...)
}
