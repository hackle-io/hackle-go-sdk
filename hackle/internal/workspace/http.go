package workspace

import (
	"encoding/json"
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	nethttp "net/http"
)

type HttpFetcher interface {
	FetchIfModified() (Workspace, bool, error)
}

func NewHttpFetcher(sdkUrl string, sdk model.Sdk, httpClient http.Client) HttpFetcher {
	return &httpFetcher{
		url:          sdkUrl + "/api/v2/workspaces/" + sdk.Key + "/config",
		httpClient:   httpClient,
		lastModified: nil,
	}
}

type httpFetcher struct {
	url          string
	httpClient   http.Client
	lastModified *string
}

func (f *httpFetcher) FetchIfModified() (Workspace, bool, error) {
	req, err := f.createRequest()
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch workspace: %v", err)
	}
	res, err := f.httpClient.Execute(req)
	if err != nil {
		return nil, false, fmt.Errorf("failed to fetch workspace: %v", err)
	}
	return f.handleResponse(res)
}

func (f *httpFetcher) createRequest() (*nethttp.Request, error) {
	req, err := nethttp.NewRequest(nethttp.MethodGet, f.url, nil)
	if err != nil {
		return nil, err
	}
	if f.lastModified != nil {
		req.Header.Set("If-Modified-Since", *f.lastModified)
	}
	return req, nil
}

func (f *httpFetcher) handleResponse(res *nethttp.Response) (Workspace, bool, error) {
	defer func() {
		body := res.Body
		err := body.Close()
		if err != nil {
			logger.Warn("failed to close response body: %v", err)
		}
	}()

	if http.IsNotModified(res) {
		logger.Info("Not-Modified")
		return nil, false, nil
	}
	if !http.IsSuccessful(res) {
		return nil, false, fmt.Errorf("http status code: %d", res.StatusCode)
	}

	lastModified := res.Header.Get("Last-Modified")
	f.lastModified = &lastModified

	var dto WorkspaceDTO
	err := json.NewDecoder(res.Body).Decode(&dto)
	if err != nil {
		return nil, false, err
	}

	logger.Info("OK")
	return NewFrom(dto), true, nil
}
