package monitoring

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics/flush"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"time"
)

type MetricRegistry struct {
	*flush.BaseMetricRegistry
	url        string
	httpClient http.Client
}

func NewMetricRegistry(monitoringUrl string, scheduler schedule.Scheduler, pushInterval time.Duration, httpClient http.Client) *MetricRegistry {
	registry := &MetricRegistry{
		url:        monitoringUrl + "/metrics",
		httpClient: httpClient,
	}
	registry.BaseMetricRegistry = flush.NewBaseMetricRegistry(registry, scheduler, pushInterval)
	return registry
}

func (r *MetricRegistry) Flush(ms []metrics.Metric) {
	for _, m := range chunk(filter(ms), 500) {
		r.dispatch(m)
	}
}

func (r *MetricRegistry) dispatch(ms []metrics.Metric) {
	batch := make([]map[string]interface{}, len(ms))
	for i, m := range ms {
		batch[i] = metric(m)
	}
	dto := map[string][]map[string]interface{}{
		"metrics": batch,
	}
	err := r.httpClient.PostObj(r.url, dto)
	if err != nil {
		logger.Error("Failed to flushing metrics: %v", err)
	}
}

func filter(ms []metrics.Metric) []metrics.Metric {
	filtered := make([]metrics.Metric, 0)
	for _, m := range ms {
		if isDispatchTarget(m) {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

func isDispatchTarget(metric metrics.Metric) bool {
	switch m := metric.(type) {
	case metrics.Counter:
		return m.Count() > 0
	case metrics.Timer:
		return m.Count() > 0
	default:
		return false
	}
}

func chunk(ms []metrics.Metric, size int) [][]metrics.Metric {
	var chunks [][]metrics.Metric
	for i := 0; i < len(ms); i += size {
		end := i + size
		if end > len(ms) {
			end = len(ms)
		}
		chunks = append(chunks, ms[i:end])
	}
	return chunks
}

func metric(metric metrics.Metric) map[string]interface{} {
	id := metric.ID()
	measurements := make(map[string]float64)
	for _, measurement := range metric.Measure() {
		measurements[measurement.Field.String()] = measurement.Value()
	}
	return map[string]interface{}{
		"name":         id.Name,
		"tags":         map[string]string(id.Tags),
		"type":         string(id.Type),
		"measurements": measurements,
	}
}
