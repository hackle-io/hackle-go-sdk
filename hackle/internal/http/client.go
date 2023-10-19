package http

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"net/http"
	"strconv"
	"time"
)

type Client interface {
	Execute(req *http.Request) (*http.Response, error)
}

func NewClient(sdk model.Sdk, clock clock.Clock, timeout time.Duration) Client {
	return &httpClient{
		delegate: &http.Client{
			Timeout: timeout,
		},
		sdk:   sdk,
		clock: clock,
	}
}

type delegate interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	delegate delegate
	sdk      model.Sdk
	clock    clock.Clock
}

func (c *httpClient) Execute(req *http.Request) (*http.Response, error) {
	c.decorate(req)
	return c.delegate.Do(req)
}

func (c *httpClient) decorate(req *http.Request) {
	req.Header.Set("X-HACKLE-SDK-KEY", c.sdk.Key)
	req.Header.Set("X-HACKLE-SDK-NAME", c.sdk.Name)
	req.Header.Set("X-HACKLE-SDK-VERSION", c.sdk.Version)
	req.Header.Set("X-HACKLE-SDK-TIME", strconv.FormatInt(c.clock.CurrentMillis(), 10))
}

func IsSuccessful(res *http.Response) bool {
	return res.StatusCode >= 200 && res.StatusCode < 300
}

func IsNotModified(res *http.Response) bool {
	return res.StatusCode == 304
}
