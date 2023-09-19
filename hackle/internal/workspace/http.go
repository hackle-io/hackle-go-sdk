package workspace

import "github.com/hackle-io/hackle-go-sdk/hackle/internal/http"

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
	var dto WorkspaceDTO
	err := f.httpClient.GetObj(f.url, &dto)
	if err != nil {
		return nil, err
	}
	return NewFrom(dto), nil
}
