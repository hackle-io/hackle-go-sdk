package schedule

import (
	"time"
)

type tickerScheduler struct{}

func NewTickerScheduler() Scheduler {
	return &tickerScheduler{}
}

func (s *tickerScheduler) SchedulePeriodically(period time.Duration, task func()) Job {
	job := &tickerJob{
		stop: make(chan bool),
	}

	ticker := time.NewTicker(period)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				task()
			case <-job.stop:
				return
			}
		}
	}()
	return job
}

type tickerJob struct {
	stop chan bool
}

func (j *tickerJob) Cancel() {
	close(j.stop)
}
