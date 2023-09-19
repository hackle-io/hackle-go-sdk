package workspace

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"time"
)

type PollingFetcher struct {
	httpFetcher      HttpFetcher
	pollingInterval  time.Duration
	scheduler        schedule.Scheduler
	currentWorkspace Workspace
	pollingJob       schedule.Job
}

func NewPollingFetcher(httpFetcher HttpFetcher, pollingInterval time.Duration, scheduler schedule.Scheduler) *PollingFetcher {
	return &PollingFetcher{
		httpFetcher:      httpFetcher,
		pollingInterval:  pollingInterval,
		scheduler:        scheduler,
		currentWorkspace: nil,
		pollingJob:       nil,
	}
}

func (f *PollingFetcher) Fetch() (Workspace, bool) {
	if f.currentWorkspace != nil {
		return f.currentWorkspace, true
	} else {
		return nil, false
	}
}

func (f *PollingFetcher) Start() {
	if f.pollingJob == nil {
		f.poll()
		f.pollingJob = f.scheduler.SchedulePeriodically(f.pollingInterval, f.poll)
	}
}

func (f *PollingFetcher) Close() {
	if f.pollingJob != nil {
		f.pollingJob.Cancel()
	}
}

func (f *PollingFetcher) poll() {
	ws, err := f.httpFetcher.Fetch()
	if err != nil {
		logger.Error("Failed to poll workspace: %v", err)
		return
	}
	f.currentWorkspace = ws
}
