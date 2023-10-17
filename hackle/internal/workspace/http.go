package workspace

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/monitoring/metric"
)

type HttpFetcher interface {
	Fetch() (Workspace, error)
}

func NewHttpFetcher(baseSdkUrl string, httpClient http.Client) HttpFetcher {
	return &httpFetcher{
		url:        baseSdkUrl + "/api/v2/workspaces",
		httpClient: httpClient,
	}
}

type httpFetcher struct {
	url        string
	httpClient http.Client
}

func (f *httpFetcher) Fetch() (Workspace, error) {
	sample := metrics.NewTimerSample(clock.System)
	ws, err := f.fetch()
	metric.RecordAPI(metric.OperationGetWorkspace, sample, err == nil)
	return ws, err
}

func (f *httpFetcher) fetch() (Workspace, error) {
	var dto WorkspaceDTO
	err := f.httpClient.GetObj(f.url, &dto)
	if err != nil {
		return nil, err
	}
	return NewFrom(dto), nil
}
