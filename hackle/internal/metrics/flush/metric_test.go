package flush

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBaseMetric(t *testing.T) {

	t.Run("self", func(t *testing.T) {
		metric := &mockFlushMetric{}
		sut := NewBaseMetric(metric)
		assert.Equal(t, metric, sut.Current())
		assert.Equal(t, metric, sut.Flush())
	})
}

type mockFlushMetric struct {
	*BaseFlushMetric
}

func (m *mockFlushMetric) ID() metrics.ID {
	return metrics.ID{}
}

func (m *mockFlushMetric) Measure() []metrics.Measurement {
	return []metrics.Measurement{}
}

func (m *mockFlushMetric) InitialMetric() metrics.Metric {
	return m
}
