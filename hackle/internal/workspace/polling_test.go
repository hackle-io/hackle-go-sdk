package workspace

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestPollingFetcher_Fetch(t *testing.T) {
	type fields struct {
		httpFetcher     *mockHttpFetcher
		pollingInterval time.Duration
		scheduler       schedule.Scheduler
	}
	tests := []struct {
		name  string
		given fields
		when  func(sut *PollingFetcher) (Workspace, bool)
		then  func(sut *PollingFetcher, ws Workspace, ok bool)
	}{
		{
			name: "when before poll then return nil",
			given: fields{
				httpFetcher:     &mockHttpFetcher{returns: []interface{}{}},
				pollingInterval: 10 * time.Second,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, ws Workspace, ok bool) {
				assert.Equal(t, false, ok)
			},
		},
		{
			name: "when failed to poll then return nil",
			given: fields{
				httpFetcher:     &mockHttpFetcher{returns: []interface{}{errors.New("fail")}},
				pollingInterval: 10 * time.Second,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, ws Workspace, ok bool) {
				assert.Equal(t, false, ok)
			},
		},
		{
			name: "when workspace is fetched then return workspace",
			given: fields{
				httpFetcher:     &mockHttpFetcher{returns: []interface{}{mocks.CreateWorkspace()}},
				pollingInterval: 10 * time.Second,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, ws Workspace, ok bool) {
				assert.NotNil(t, ws)
				assert.Equal(t, true, ok)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewPollingFetcher(
				tt.given.httpFetcher,
				tt.given.pollingInterval,
				tt.given.scheduler,
			)
			ws, ok := tt.when(sut)
			tt.then(sut, ws, ok)
			sut.Close()
		})
	}
}

func TestPollingFetcher_poll(t *testing.T) {
	type fields struct {
		httpFetcher     *mockHttpFetcher
		pollingInterval time.Duration
		scheduler       schedule.Scheduler
	}
	tests := []struct {
		name  string
		given fields
		when  func(sut *PollingFetcher) (Workspace, bool)
		then  func(sut *PollingFetcher, ws Workspace, ok bool)
	}{
		{
			name: "failed to poll",
			given: fields{
				httpFetcher:     &mockHttpFetcher{returns: []interface{}{errors.New("fail")}},
				pollingInterval: 10 * time.Second,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, ws Workspace, ok bool) {
				assert.Equal(t, false, ok)
			},
		},
		{
			name: "success to poll",
			given: fields{
				httpFetcher:     &mockHttpFetcher{returns: []interface{}{mocks.CreateWorkspace()}},
				pollingInterval: 10 * time.Second,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, ws Workspace, ok bool) {
				assert.NotNil(t, ws)
				assert.Equal(t, true, ok)
			},
		},
		{
			name: "workspace not modified",
			given: fields{
				httpFetcher: func() *mockHttpFetcher {
					return &mockHttpFetcher{returns: []interface{}{mocks.CreateWorkspace(), nil, nil, nil, nil}}
				}(),
				pollingInterval: 100 * time.Millisecond,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				time.Sleep(350 * time.Millisecond)
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, ws Workspace, ok bool) {
				assert.NotNil(t, ws)
				assert.Equal(t, true, ok)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewPollingFetcher(
				tt.given.httpFetcher,
				tt.given.pollingInterval,
				tt.given.scheduler,
			)
			ws, ok := tt.when(sut)
			tt.then(sut, ws, ok)
			sut.Close()
		})
	}
}

func TestPollingFetcher_Start(t *testing.T) {
	type fields struct {
		httpFetcher     *mockHttpFetcher
		pollingInterval time.Duration
		scheduler       schedule.Scheduler
	}
	tests := []struct {
		name  string
		given fields
		when  func(sut *PollingFetcher) (Workspace, bool)
		then  func(sut *PollingFetcher, given fields, ws Workspace, ok bool)
	}{
		{
			name: "poll",
			given: fields{
				httpFetcher:     &mockHttpFetcher{returns: []interface{}{mocks.CreateWorkspace()}},
				pollingInterval: 10 * time.Second,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, given fields, ws Workspace, ok bool) {
				assert.NotNil(t, ws)
				assert.Equal(t, true, ok)
			},
		},
		{
			name: "start scheduling",
			given: fields{
				httpFetcher: &mockHttpFetcher{returns: []interface{}{
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
				}},
				pollingInterval: 100 * time.Millisecond,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				sut.Start()
				time.Sleep(550 * time.Millisecond)
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, given fields, ws Workspace, ok bool) {
				assert.NotNil(t, ws)
				assert.Equal(t, true, ok)
				assert.Equal(t, 6, given.httpFetcher.Count())
			},
		},
		{
			name: "start once",
			given: fields{
				httpFetcher: &mockHttpFetcher{returns: []interface{}{
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
					mocks.CreateWorkspace(),
				}},
				pollingInterval: 100 * time.Millisecond,
				scheduler:       schedule.NewTickerScheduler(),
			},
			when: func(sut *PollingFetcher) (Workspace, bool) {
				for i := 0; i < 10; i++ {
					sut.Start()
				}
				time.Sleep(550 * time.Millisecond)
				return sut.Fetch()
			},
			then: func(sut *PollingFetcher, given fields, ws Workspace, ok bool) {
				assert.NotNil(t, ws)
				assert.Equal(t, true, ok)
				assert.Equal(t, 6, given.httpFetcher.Count())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := NewPollingFetcher(
				tt.given.httpFetcher,
				tt.given.pollingInterval,
				tt.given.scheduler,
			)
			ws, ok := tt.when(sut)
			tt.then(sut, tt.given, ws, ok)
			sut.Close()
		})
	}
}

type mockHttpFetcher struct {
	returns []interface{}
	count   int
	mu      sync.Mutex
}

func (m *mockHttpFetcher) FetchIfModified() (Workspace, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ret := m.returns[m.count]
	m.count++
	switch r := ret.(type) {
	case Workspace:
		return r, true, nil
	case error:
		return nil, false, r
	default:
		return nil, false, nil
	}
}

func (m *mockHttpFetcher) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.count
}
