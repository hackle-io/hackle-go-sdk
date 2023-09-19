package schedule

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTickerScheduler(t *testing.T) {

	t.Run("schedule", func(t *testing.T) {
		count := 0

		task := func() {
			count++
		}

		scheduler := NewTickerScheduler()
		job := scheduler.SchedulePeriodically(300*time.Millisecond, 500*time.Millisecond, task)
		time.Sleep(1100 * time.Millisecond)
		job.Cancel()
		assert.Equal(t, 2, count)
	})

	t.Run("cancel during delay", func(t *testing.T) {
		count := 0

		task := func() {
			count++
		}

		scheduler := NewTickerScheduler()
		job := scheduler.SchedulePeriodically(1000*time.Millisecond, 100*time.Millisecond, task)
		time.Sleep(500 * time.Millisecond)
		job.Cancel()
		assert.Equal(t, 0, count)
	})

	t.Run("long task", func(t *testing.T) {
		count := 0

		longTask := func() {
			time.Sleep(400 * time.Millisecond)
			count++
		}
		scheduler := NewTickerScheduler()
		job := scheduler.SchedulePeriodically(0*time.Millisecond, 100*time.Millisecond, longTask)
		time.Sleep(1000 * time.Millisecond)
		job.Cancel()
		assert.Equal(t, 2, count)
	})
}
