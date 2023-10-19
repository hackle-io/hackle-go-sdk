package workspace

import (
	"bytes"
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	nethttp "net/http"
	"testing"
)

func TestNewHttpFetcher(t *testing.T) {
	fetcher := NewHttpFetcher("localhost", model.Sdk{Key: "sdk_key"}, &mockHttpClient{})
	assert.IsType(t, &httpFetcher{}, fetcher)
	assert.Equal(t, "localhost/api/v2/workspaces/sdk_key/config", fetcher.(*httpFetcher).url)
}

func Test_httpFetcher_FetchIfModified(t *testing.T) {
	type fields struct {
		url          string
		httpClient   *mockHttpClient
		lastModified *string
	}
	tests := []struct {
		name      string
		fields    fields
		assertion func(*httpFetcher, fields, Workspace, bool, error)
	}{
		{
			name: "when failed to create request then return error",
			fields: fields{
				url:          ":",
				httpClient:   &mockHttpClient{},
				lastModified: nil,
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				assert.Equal(t, nil, ws)
				assert.Equal(t, false, ok)
				assert.Contains(t, err.Error(), "failed to fetch workspace")
			},
		},
		{
			name: "when failed to http execute then return error",
			fields: fields{
				url:          "localhost",
				httpClient:   &mockHttpClient{err: errors.New("http fail")},
				lastModified: nil,
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				assert.Equal(t, nil, ws)
				assert.Equal(t, false, ok)
				assert.Contains(t, err.Error(), "failed to fetch workspace: http fail")
			},
		},
		{
			name: "when workspace not modified then return nil",
			fields: fields{
				url: "localhost",
				httpClient: &mockHttpClient{
					res: &nethttp.Response{
						StatusCode: 304,
						Body:       ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
					},
				},
				lastModified: nil,
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				assert.Equal(t, nil, ws)
				assert.Equal(t, false, ok)
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "when http not successful then return error",
			fields: fields{
				url: "localhost",
				httpClient: &mockHttpClient{
					res: &nethttp.Response{
						StatusCode: 500,
						Body:       ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
					},
				},
				lastModified: nil,
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				assert.Equal(t, nil, ws)
				assert.Equal(t, false, ok)
				assert.Contains(t, err.Error(), "http status code: 500")
			},
		},
		{
			name: "when success to get workspace then return workspace",
			fields: fields{
				url: "localhost",
				httpClient: func() *mockHttpClient {
					body, _ := ioutil.ReadFile("../../../testdata/workspace_config.json")
					return &mockHttpClient{
						res: &nethttp.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewReader(body)),
						},
					}
				}(),
				lastModified: nil,
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				_, a := ws.GetExperiment(5)
				assert.Equal(t, true, a)
			},
		},
		{
			name: "when response Last-Modified header then set header",
			fields: fields{
				url: "localhost",
				httpClient: func() *mockHttpClient {
					body, _ := ioutil.ReadFile("../../../testdata/workspace_config.json")
					return &mockHttpClient{
						res: &nethttp.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewReader(body)),
							Header: nethttp.Header{
								"Last-Modified": {"LAST_MODIFIED_HEADER_VALUE"},
							},
						},
					}
				}(),
				lastModified: nil,
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				assert.Equal(t, "LAST_MODIFIED_HEADER_VALUE", *sut.lastModified)
			},
		},
		{
			name: "when Last-Modified header is exist then execute with header",
			fields: fields{
				url: "localhost",
				httpClient: func() *mockHttpClient {
					body, _ := ioutil.ReadFile("../../../testdata/workspace_config.json")
					return &mockHttpClient{
						res: &nethttp.Response{
							StatusCode: 200,
							Body:       ioutil.NopCloser(bytes.NewReader(body)),
						},
					}
				}(),
				lastModified: ref.String("LAST_MODIFIED_HEADER_VALUE"),
			},
			assertion: func(sut *httpFetcher, fields fields, ws Workspace, ok bool, err error) {
				assert.Equal(t, "LAST_MODIFIED_HEADER_VALUE", fields.httpClient.req.Header.Get("If-Modified-Since"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sut := &httpFetcher{
				url:          tt.fields.url,
				httpClient:   tt.fields.httpClient,
				lastModified: tt.fields.lastModified,
			}
			ws, ok, err := sut.FetchIfModified()
			tt.assertion(sut, tt.fields, ws, ok, err)
		})
	}
}

type mockHttpClient struct {
	req *nethttp.Request
	res *nethttp.Response
	err error
}

func (m *mockHttpClient) Execute(req *nethttp.Request) (*nethttp.Response, error) {
	m.req = req
	if m.err != nil {
		return nil, m.err
	}
	return m.res, nil
}
