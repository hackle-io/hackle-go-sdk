package workspace

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPollingFetcher_Fetch(t *testing.T) {

	t.Run("without start", func(t *testing.T) {
		sut := NewPollingFetcher(&mockHttpFetcher{}, 10*time.Second, schedule.NewTickerScheduler())
		ws, ok := sut.Fetch()
		assert.Equal(t, nil, ws)
		assert.Equal(t, false, ok)
	})

	t.Run("poll fail", func(t *testing.T) {
		sut := NewPollingFetcher(&mockHttpFetcher{returns: errors.New("poll fail")}, 100*time.Millisecond, schedule.NewTickerScheduler())
		sut.Start()
		time.Sleep(500)
		ws, ok := sut.Fetch()
		assert.Equal(t, nil, ws)
		assert.Equal(t, false, ok)
	})
}

func TestPollingFetcher_Start(t *testing.T) {
	t.Run("poll", func(t *testing.T) {
		sut := NewPollingFetcher(&mockHttpFetcher{returns: mocks.CreateWorkspace()}, 10*time.Second, schedule.NewTickerScheduler())
		sut.Start()

		ws, ok := sut.Fetch()
		assert.NotNil(t, ws)
		assert.Equal(t, true, ok)
	})

	t.Run("once", func(t *testing.T) {
		scheduler := &mockScheduler{}
		sut := NewPollingFetcher(&mockHttpFetcher{returns: mocks.CreateWorkspace()}, 10*time.Second, scheduler)

		for i := 0; i < 10; i++ {
			sut.Start()
		}
		assert.Equal(t, 1, scheduler.jobCount())
	})

	t.Run("long task polling", func(t *testing.T) {
		httpFetcher := &mockHttpFetcher{returns: mocks.CreateWorkspace(), delay: 250 * time.Millisecond}
		sut := NewPollingFetcher(httpFetcher, 100*time.Millisecond, schedule.NewTickerScheduler())

		sut.Start()
		time.Sleep(400 * time.Millisecond)

		assert.Equal(t, 3, httpFetcher.count) // 0, 250, 350 fetch / 100, 200 ignored
	})
}

func TestPollingFetcher_Close(t *testing.T) {
	t.Run("cancel job", func(t *testing.T) {
		scheduler := &mockScheduler{}
		sut := NewPollingFetcher(&mockHttpFetcher{returns: mocks.CreateWorkspace()}, 10*time.Second, scheduler)

		sut.Start()
		job := scheduler.jobs[0]
		assert.Equal(t, false, job.canceled)

		sut.Close()
		assert.Equal(t, true, job.canceled)
	})
}

type mockHttpFetcher struct {
	delay   time.Duration
	returns interface{}
	count   int
}

func (m *mockHttpFetcher) Fetch() (Workspace, error) {
	m.count++
	time.Sleep(m.delay)
	switch r := m.returns.(type) {
	case Workspace:
		return r, nil
	case error:
		return nil, r
	default:
		return nil, nil
	}
}

type mockScheduler struct {
	jobs []*mockJob
}

func (m *mockScheduler) SchedulePeriodically(delay time.Duration, period time.Duration, task func()) schedule.Job {
	job := &mockJob{}
	m.jobs = append(m.jobs, job)
	return job
}

func (m *mockScheduler) jobCount() int {
	return len(m.jobs)
}

type mockJob struct {
	canceled bool
}

func (m *mockJob) Cancel() {
	m.canceled = true
}
