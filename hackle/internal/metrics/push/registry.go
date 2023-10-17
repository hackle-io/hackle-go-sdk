package push

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"sync"
	"time"
)

type MetricRegistry interface {
	Publish()
	Start()
	Stop()
}

type BaseMetricRegistry struct {
	MetricRegistry
	scheduler     schedule.Scheduler
	pushInterval  time.Duration
	publishingJob schedule.Job
	mu            sync.Mutex
}

func NewBaseMetricRegistry(self MetricRegistry, scheduler schedule.Scheduler, pushInterval time.Duration) *BaseMetricRegistry {
	registry := &BaseMetricRegistry{
		MetricRegistry: self,
		scheduler:      scheduler,
		pushInterval:   pushInterval,
	}
	return registry
}

func (r *BaseMetricRegistry) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.publishingJob != nil {
		return
	}
	r.publishingJob = r.scheduler.SchedulePeriodically(r.pushInterval, r.Publish)
	logger.Info("%T started. Publish every %v", r.MetricRegistry, r.pushInterval)
}

func (r *BaseMetricRegistry) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.publishingJob == nil {
		return
	}
	r.publishingJob.Cancel()
	r.publishingJob = nil
	r.Publish()
}
