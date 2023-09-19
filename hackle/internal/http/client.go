package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"net/http"
	"strconv"
	"time"
)

type Client interface {
	GetObj(url string, out interface{}) error
	PostObj(url string, body interface{}) error
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

func (c *httpClient) GetObj(url string, out interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	c.setHeaders(req)

	res, err := c.delegate.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		body := res.Body
		e := body.Close()
		if e != nil {
			logger.Warn("failed to close response body: %v", e)
		}
	}()

	if !isSuccessful(res.StatusCode) {
		return fmt.Errorf("http status code: %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(out)
}

func (c *httpClient) PostObj(url string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	c.setHeaders(req)

	res, err := c.delegate.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		e := res.Body.Close()
		if e != nil {
			logger.Warn("failed to close response body: %v", e)
		}
	}()

	if !isSuccessful(res.StatusCode) {
		return fmt.Errorf("http status code: %d", res.StatusCode)
	}

	return nil
}

func (c *httpClient) setHeaders(request *http.Request) {
	request.Header.Set("X-HACKLE-SDK-KEY", c.sdk.Key)
	request.Header.Set("X-HACKLE-SDK-NAME", c.sdk.Name)
	request.Header.Set("X-HACKLE-SDK-VERSION", c.sdk.Version)
	request.Header.Set("X-HACKLE-SDK-TIME", strconv.FormatInt(c.clock.CurrentMillis(), 10))
	request.Header.Set("Content-Type", "application/json")
}

func isSuccessful(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
