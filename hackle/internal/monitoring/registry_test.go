package monitoring

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewMetricRegistry(t *testing.T) {

	registry := NewMetricRegistry(
		"localhost",
		schedule.NewTickerScheduler(),
		10*time.Second,
		http.NewClient(model.Sdk{}, clock.System, 10*time.Second),
	)

	assert.Equal(t, "localhost/metrics", registry.url)
	assert.Equal(t, registry, registry.BaseMetricRegistry.MetricRegistry)
}

func TestMetricRegistry_Flush(t *testing.T) {

	t.Run("when metric is empty then do not dispatch", func(t *testing.T) {
		httpClient := &mockHttpClient{}
		sut := NewMetricRegistry(
			"localhost",
			schedule.NewTickerScheduler(),
			10*time.Second,
			httpClient,
		)

		sut.Flush(make([]metrics.Metric, 0))

		assert.Equal(t, 0, len(httpClient.Posts()))
	})

	t.Run("when metric is not dispatch target then do not dispatch", func(t *testing.T) {
		httpClient := &mockHttpClient{}
		sut := NewMetricRegistry(
			"localhost",
			schedule.NewTickerScheduler(),
			10*time.Second,
			httpClient,
		)

		registry := metrics.NewCumulativeRegistry()
		counter := registry.Counter("counter", metrics.Tags{})
		timer := registry.Timer("timer", metrics.Tags{})

		sut.Flush([]metrics.Metric{counter, timer})

		assert.Equal(t, 0, len(httpClient.Posts()))
	})

	t.Run("dispatch", func(t *testing.T) {
		httpClient := &mockHttpClient{}
		sut := NewMetricRegistry(
			"localhost",
			schedule.NewTickerScheduler(),
			10*time.Second,
			httpClient,
		)

		registry := metrics.NewCumulativeRegistry()
		counter := registry.Counter("counter", metrics.Tags{"a": "1"})
		counter.Increment(1)
		timer := registry.Timer("timer", metrics.Tags{"b": "2"})
		timer.Record(42 * time.Millisecond)

		sut.Flush([]metrics.Metric{counter, timer})

		assert.Equal(t, 1, len(httpClient.Posts()))
		assert.Equal(t, map[string][]map[string]interface{}{
			"metrics": {
				{
					"name": "counter",
					"tags": map[string]string{"a": "1"},
					"type": "COUNTER",
					"measurements": map[string]float64{
						"count": 1.0,
					},
				},
				{
					"name": "timer",
					"tags": map[string]string{"b": "2"},
					"type": "TIMER",
					"measurements": map[string]float64{
						"count": 1.0,
						"total": 42.0,
						"max":   42.0,
						"mean":  42.0,
					},
				},
			},
		}, httpClient.Posts()[0])
	})

	t.Run("dispatch failed", func(t *testing.T) {
		httpClient := &mockHttpClient{err: errors.New("fail")}
		sut := NewMetricRegistry(
			"localhost",
			schedule.NewTickerScheduler(),
			10*time.Second,
			httpClient,
		)

		registry := metrics.NewCumulativeRegistry()
		counter := registry.Counter("counter", metrics.Tags{"a": "1"})
		counter.Increment(1)
		timer := registry.Timer("timer", metrics.Tags{"b": "2"})
		timer.Record(42 * time.Millisecond)

		sut.Flush([]metrics.Metric{counter, timer})
	})
}

type mockHttpClient struct {
	posts []interface{}
	err   error
}

func (m *mockHttpClient) GetObj(url string, out interface{}) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockHttpClient) PostObj(url string, body interface{}) error {
	if m.err != nil {
		return m.err
	}
	m.posts = append(m.posts, body)
	return nil
}

func (m *mockHttpClient) Posts() []interface{} {
	return m.posts
}
