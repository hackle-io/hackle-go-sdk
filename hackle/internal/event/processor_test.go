package event

import (
	"github.com/google/uuid"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

type fields struct {
	queue         chan message
	dispatcher    *mockDispatcher
	consumingWait *sync.WaitGroup
}

func eventProcessor(capacity int, dispatchSize int, flushInterval time.Duration) (*processor, *fields) {
	f := &fields{
		queue:         make(chan message, capacity),
		dispatcher:    &mockDispatcher{},
		consumingWait: &sync.WaitGroup{},
	}
	return &processor{
		queue:          f.queue,
		dispatcher:     f.dispatcher,
		dispatchSize:   dispatchSize,
		flushScheduler: schedule.NewTickerScheduler(),
		flushInterval:  flushInterval,
		consumingWait:  f.consumingWait,
		flushingJob:    nil,
		isStarted:      false,
	}, f
}

func TestNewProcessor(t *testing.T) {
	p := NewProcessor(42, &mockDispatcher{}, 320, &mockScheduler{}, 100*time.Millisecond)
	assert.IsType(t, &processor{}, p)
}

func TestProcessor_Process(t *testing.T) {

	t.Run("put event message in the queue", func(t *testing.T) {
		sut, f := eventProcessor(10, 100, 10*time.Second)

		assert.Equal(t, 0, len(f.queue))

		event := baseUserEvent{insertID: "42"}
		sut.Process(event)
		assert.Equal(t, 1, len(f.queue))

		msg := <-f.queue
		assert.IsType(t, eventMessage{}, msg)
		assert.Equal(t, "42", msg.(eventMessage).event.InsertID())
	})

	t.Run("when queue is full then ignore", func(t *testing.T) {
		sut, f := eventProcessor(10, 100, 10*time.Second)

		assert.Equal(t, 0, len(f.queue))
		for i := 0; i < 10; i++ {
			sut.Process(baseUserEvent{})
		}
		assert.Equal(t, 10, len(f.queue))

		sut.Process(baseUserEvent{})
		assert.Equal(t, 10, len(f.queue))
	})

	t.Run("when dispatch size not reached then do not dispatch", func(t *testing.T) {
		sut, f := eventProcessor(100, 2, 10*time.Second)
		sut.Start()

		sut.Process(baseUserEvent{})
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 0, f.dispatcher.DispatchCount())
	})

	t.Run("when dispatch size reached then dispatch event", func(t *testing.T) {
		sut, f := eventProcessor(100, 2, 10*time.Second)
		sut.Start()

		sut.Process(baseUserEvent{insertID: "1"})
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 0, f.dispatcher.DispatchCount())

		sut.Process(baseUserEvent{insertID: "2"})
		time.Sleep(100 * time.Millisecond)
		assert.Equal(t, 1, f.dispatcher.DispatchCount())
	})

	t.Run("when flush interval reached then dispatch events", func(t *testing.T) {
		sut, f := eventProcessor(100, 1000, 500*time.Millisecond)
		sut.Start()

		sut.Process(baseUserEvent{insertID: "1"})
		sut.Process(baseUserEvent{insertID: "2"})
		sut.Process(baseUserEvent{insertID: "3"})
		sut.Process(baseUserEvent{insertID: "4"})
		sut.Process(baseUserEvent{insertID: "5"})
		time.Sleep(700 * time.Millisecond)

		assert.Equal(t, 1, f.dispatcher.DispatchCount())
	})

	t.Run("when event is empty then not dispatch", func(t *testing.T) {
		sut, f := eventProcessor(100, 1000, 100*time.Millisecond)
		sut.Start()

		time.Sleep(1 * time.Second)

		assert.Equal(t, 0, f.dispatcher.DispatchCount())
	})

	t.Run("concurrency", func(t *testing.T) {

		sut, f := eventProcessor(16*10000, 100, 10*time.Millisecond)
		sut.Start()

		task := func() {
			for i := 0; i < 10000; i++ {
				sut.Process(baseUserEvent{insertID: uuid.NewString()})
			}
		}

		wg := sync.WaitGroup{}
		for i := 0; i < 16; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				task()
			}()
		}
		wg.Wait()
		sut.Close()

		time.Sleep(100 * time.Millisecond)
		results := map[string]bool{}
		dispatched := f.dispatcher.DispatchedEvents()
		for _, userEvent := range dispatched {
			results[userEvent.InsertID()] = true
		}
		assert.Equal(t, 16*10000, f.dispatcher.EventCount())
		assert.Equal(t, 16*10000, len(results))
	})
}

func TestProcessor_Start(t *testing.T) {
	t.Run("start once", func(t *testing.T) {
		scheduler := &mockScheduler{}
		sut, _ := eventProcessor(16*10000, 10, 10*time.Second)
		sut.flushScheduler = scheduler

		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sut.Start()
			}()
		}

		wg.Wait()

		assert.Equal(t, 1, scheduler.jobCount())
	})
}

func TestProcessor_Close(t *testing.T) {
	t.Run("cancel flushing", func(t *testing.T) {
		scheduler := &mockScheduler{}
		sut, _ := eventProcessor(16*10000, 10, 10*time.Second)
		sut.flushScheduler = scheduler

		sut.Start()
		assert.Equal(t, false, scheduler.jobs[0].canceled)

		sut.Close()
		assert.Equal(t, true, scheduler.jobs[0].canceled)
	})

	t.Run("shutdown consuming", func(t *testing.T) {
		sut, f := eventProcessor(16*10000, 10, 10*time.Second)
		sut.Start()

		sut.Process(baseUserEvent{})
		sut.Close()
		time.Sleep(100 * time.Millisecond)

		assert.Equal(t, 1, f.dispatcher.DispatchCount())
	})

	t.Run("close dispatcher", func(t *testing.T) {
		sut, f := eventProcessor(16*10000, 10, 10*time.Second)

		sut.Start()
		assert.Equal(t, false, f.dispatcher.closed)

		sut.Process(baseUserEvent{})
		sut.Close()
		assert.Equal(t, true, f.dispatcher.closed)
	})

	t.Run("not started", func(t *testing.T) {
		sut, f := eventProcessor(16*10000, 10, 10*time.Second)
		sut.Close()
		assert.Equal(t, false, f.dispatcher.closed)
	})
}

type mockDispatcher struct {
	dispatched    []UserEvent
	dispatchCount int
	eventCount    int
	closed        bool
	mu            sync.Mutex
	wg            sync.WaitGroup
}

func (m *mockDispatcher) Dispatch(userEvents []UserEvent) {
	go func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		for _, userEvent := range userEvents {
			m.dispatched = append(m.dispatched, userEvent)
		}
		m.dispatchCount++
		m.eventCount = m.eventCount + len(userEvents)
	}()
}

func (m *mockDispatcher) Close() {
	m.closed = true
}

func (m *mockDispatcher) DispatchCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dispatchCount
}

func (m *mockDispatcher) EventCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.eventCount
}

func (m *mockDispatcher) DispatchedEvents() []UserEvent {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.dispatched
}

type mockScheduler struct {
	jobs []*mockJob
}

func (m *mockScheduler) SchedulePeriodically(period time.Duration, task func()) schedule.Job {
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
