package http

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(model.Sdk{}, clock.System, 10*time.Second)
	assert.IsType(t, &httpClient{}, client)

	hc := client.(*httpClient)
	assert.IsType(t, &http.Client{}, hc.delegate)
}

func TestHttpClient_Execute(t *testing.T) {
	res := &http.Response{StatusCode: 200}
	delegate := &mockHttpClient{returns: res}
	sut := &httpClient{
		delegate: delegate,
		sdk:      model.Sdk{Key: "sdk_key", Name: "test-sdk", Version: "test-version"},
		clock:    clock.Fixed(42),
	}

	req, _ := http.NewRequest(http.MethodGet, "localhost", nil)
	actual, _ := sut.Execute(req)

	assert.Same(t, res, actual)
	assert.Same(t, req, delegate.req)
	h := http.Header{}
	h.Set("X-HACKLE-SDK-KEY", "sdk_key")
	h.Set("X-HACKLE-SDK-NAME", "test-sdk")
	h.Set("X-HACKLE-SDK-VERSION", "test-version")
	h.Set("X-HACKLE-SDK-TIME", "42")
	assert.Equal(t, h, delegate.req.Header)
}

func TestIsSuccessful(t *testing.T) {
	assert.Equal(t, true, IsSuccessful(&http.Response{StatusCode: 200}))
	assert.Equal(t, true, IsSuccessful(&http.Response{StatusCode: 299}))
	assert.Equal(t, false, IsSuccessful(&http.Response{StatusCode: 199}))
	assert.Equal(t, false, IsSuccessful(&http.Response{StatusCode: 300}))
	assert.Equal(t, false, IsSuccessful(&http.Response{StatusCode: 400}))
	assert.Equal(t, false, IsSuccessful(&http.Response{StatusCode: 500}))
}

func TestIsNotModified(t *testing.T) {
	assert.Equal(t, true, IsNotModified(&http.Response{StatusCode: 304}))
	assert.Equal(t, false, IsNotModified(&http.Response{StatusCode: 200}))
	assert.Equal(t, false, IsNotModified(&http.Response{StatusCode: 500}))
}

type mockHttpClient struct {
	req     *http.Request
	returns interface{}
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.req = req
	switch r := m.returns.(type) {
	case *http.Response:
		return r, nil
	case error:
		return nil, r
	default:
		return nil, nil
	}
}
