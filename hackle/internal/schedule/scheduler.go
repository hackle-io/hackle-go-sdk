package schedule

import "time"

type Scheduler interface {
	SchedulePeriodically(delay time.Duration, period time.Duration, task func()) Job
}

type Job interface {
	Cancel()
}
