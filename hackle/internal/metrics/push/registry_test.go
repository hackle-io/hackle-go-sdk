package push

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewBaseMetricRegistry(t *testing.T) {
	registry := newMockBaseMetricRegistry(schedule.NewTickerScheduler(), 1*time.Hour)

	assert.Equal(t, registry.MetricRegistry, registry)
}

func TestBaseMetricRegistry_Start(t *testing.T) {

	t.Run("schedule publish", func(t *testing.T) {
		registry := newMockBaseMetricRegistry(schedule.NewTickerScheduler(), 200*time.Millisecond)
		assert.Equal(t, nil, registry.publishingJob)

		registry.Start()
		assert.NotNil(t, registry.publishingJob)
		time.Sleep(500 * time.Millisecond)

		assert.Equal(t, 2, registry.PublishCount())
	})

	t.Run("start once", func(t *testing.T) {
		scheduler := &mockScheduler{}
		registry := newMockBaseMetricRegistry(scheduler, 200*time.Millisecond)
		for i := 0; i < 100; i++ {
			registry.Start()
		}

		assert.Equal(t, 1, len(scheduler.jobs))
	})
}

func TestBaseMetricRegistry_Stop(t *testing.T) {

	t.Run("cancel publishing job", func(t *testing.T) {
		scheduler := &mockScheduler{}
		registry := newMockBaseMetricRegistry(scheduler, 200*time.Millisecond)

		registry.Start()
		assert.Equal(t, false, scheduler.jobs[0].cancelled)

		registry.Stop()
		assert.Equal(t, true, scheduler.jobs[0].cancelled)
		assert.Equal(t, nil, registry.publishingJob)
	})

	t.Run("publish", func(t *testing.T) {
		scheduler := &mockScheduler{}
		registry := newMockBaseMetricRegistry(scheduler, 1*time.Hour)

		registry.Start()
		assert.Equal(t, 0, registry.PublishCount())

		registry.Stop()
		assert.Equal(t, 1, registry.PublishCount())
	})

	t.Run("not started", func(t *testing.T) {
		scheduler := &mockScheduler{}
		registry := newMockBaseMetricRegistry(scheduler, 1*time.Hour)

		registry.Stop()

		assert.Equal(t, 0, registry.PublishCount())
	})
}

func newMockBaseMetricRegistry(scheduler schedule.Scheduler, pushInterval time.Duration) *mockBaseMetricRegistry {
	registry := &mockBaseMetricRegistry{}
	registry.BaseMetricRegistry = NewBaseMetricRegistry(registry, scheduler, pushInterval)
	return registry
}

type mockBaseMetricRegistry struct {
	*BaseMetricRegistry
	publishCount int
	mu           sync.Mutex
}

func (m *mockBaseMetricRegistry) Publish() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publishCount++
}

func (m *mockBaseMetricRegistry) PublishCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.publishCount
}

type mockScheduler struct {
	jobs []*mockJob
}

func (m *mockScheduler) SchedulePeriodically(period time.Duration, task func()) schedule.Job {
	job := &mockJob{}
	m.jobs = append(m.jobs, job)
	return job
}

type mockJob struct {
	cancelled bool
}

func (m *mockJob) Cancel() {
	m.cancelled = true
}
