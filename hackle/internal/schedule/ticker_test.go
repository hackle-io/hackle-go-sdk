package schedule

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestTickerScheduler(t *testing.T) {

	t.Run("schedule", func(t *testing.T) {

		c := &counter{}

		scheduler := NewTickerScheduler()
		job := scheduler.SchedulePeriodically(500*time.Millisecond, func() {
			c.increment()
		})
		time.Sleep(1200 * time.Millisecond)
		job.Cancel()
		assert.Equal(t, 2, c.get())
	})

	t.Run("cancel during delay", func(t *testing.T) {
		c := &counter{}
		scheduler := NewTickerScheduler()
		job := scheduler.SchedulePeriodically(500*time.Millisecond, func() {
			c.increment()
		})
		time.Sleep(100 * time.Millisecond)
		job.Cancel()
		assert.Equal(t, 0, c.get())
	})

	t.Run("long task", func(t *testing.T) {
		c := &counter{}
		scheduler := NewTickerScheduler()
		job := scheduler.SchedulePeriodically(100*time.Millisecond, func() {
			time.Sleep(400 * time.Millisecond)
			c.increment()
		})
		time.Sleep(1000 * time.Millisecond)
		job.Cancel()
		assert.Equal(t, 2, c.get())
	})
}

type counter struct {
	count int
	mu    sync.Mutex
}

func (c *counter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *counter) get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}
