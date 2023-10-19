package workspace

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"sync"
	"time"
)

type PollingFetcher struct {
	httpFetcher      HttpFetcher
	pollingInterval  time.Duration
	scheduler        schedule.Scheduler
	currentWorkspace Workspace
	pollingJob       schedule.Job
	mu               sync.Mutex
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
	f.mu.Lock()
	defer f.mu.Unlock()
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
	ws, ok, err := f.httpFetcher.FetchIfModified()
	if err != nil {
		logger.Error("Failed to poll workspace: %v", err)
		return
	}
	if ok {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.currentWorkspace = ws
	}
}
