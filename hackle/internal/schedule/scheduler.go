package schedule

import "time"

type Scheduler interface {
	SchedulePeriodically(period time.Duration, task func()) Job
}

type Job interface {
	Cancel()
}
