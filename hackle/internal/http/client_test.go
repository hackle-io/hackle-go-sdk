package http

import (
	"errors"
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

func Test_httpClient_GetObj(t *testing.T) {
	type fields struct {
		delegate *mockHttpClient
		sdk      model.Sdk
		clock    clock.Clock
	}
	type args struct {
		url string
		out interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion func(fields fields, args args, err error)
	}{
		{
			name: "failed to create request",
			fields: fields{
				delegate: &mockHttpClient{returns: nil},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url: ":",
				out: nil,
			},
			assertion: func(fields fields, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "failed to request",
			fields: fields{
				delegate: &mockHttpClient{returns: errors.New("do fail")},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url: "http://localhost",
				out: nil,
			},
			assertion: func(fields fields, args args, err error) {
				assert.Equal(t, errors.New("do fail"), err)
			},
		},
		{
			name: "not successful",
			fields: fields{
				delegate: &mockHttpClient{returns: &http.Response{StatusCode: 500, Body: &mockBody{}}},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url: "http://localhost",
				out: nil,
			},
			assertion: func(fields fields, args args, err error) {
				assert.Equal(t, errors.New("http status code: 500"), err)
			},
		},
		{
			name: "successful",
			fields: fields{
				delegate: &mockHttpClient{returns: &http.Response{StatusCode: 200, Body: &mockBody{bytes: []byte(`{"a":"1","b":"2"}`)}}},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url: "http://localhost",
				out: &data{},
			},
			assertion: func(fields fields, args args, err error) {
				assert.Equal(t, data{A: "1", B: "2"}, *args.out.(*data))
				assert.Equal(t, nil, err)
			},
		},
		{
			name: "add header",
			fields: fields{
				delegate: &mockHttpClient{returns: &http.Response{StatusCode: 200, Body: &mockBody{bytes: []byte(`{"a":"1","b":"2"}`)}}},
				sdk:      model.Sdk{Key: "sdk_key", Name: "test-sdk", Version: "test-version"},
				clock:    clock.Fixed(42),
			},
			args: args{
				url: "http://localhost",
				out: &data{},
			},
			assertion: func(fields fields, args args, err error) {
				h := http.Header{}
				h.Set("X-HACKLE-SDK-KEY", fields.sdk.Key)
				h.Set("X-HACKLE-SDK-NAME", fields.sdk.Name)
				h.Set("X-HACKLE-SDK-VERSION", fields.sdk.Version)
				h.Set("X-HACKLE-SDK-TIME", "42")
				h.Set("Content-Type", "application/json")
				assert.Equal(t, h, fields.delegate.req.Header)
				assert.Equal(t, data{A: "1", B: "2"}, *args.out.(*data))
				assert.Equal(t, nil, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpClient{
				delegate: tt.fields.delegate,
				sdk:      tt.fields.sdk,
				clock:    tt.fields.clock,
			}
			err := c.GetObj(tt.args.url, tt.args.out)
			tt.assertion(tt.fields, tt.args, err)
		})
	}
}

func Test_httpClient_PostObj(t *testing.T) {
	type fields struct {
		delegate *mockHttpClient
		sdk      model.Sdk
		clock    clock.Clock
	}
	type args struct {
		url  string
		body interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion func(fields fields, args args, err error)
	}{
		{
			name: "marshal fail",
			fields: fields{
				delegate: &mockHttpClient{returns: nil},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url:  "localhost",
				body: make(chan int),
			},
			assertion: func(fields fields, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "create request fail",
			fields: fields{
				delegate: &mockHttpClient{returns: nil},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url:  ":",
				body: data{A: "1", B: "2"},
			},
			assertion: func(fields fields, args args, err error) {
				assert.NotNil(t, err)
			},
		},
		{
			name: "post fail",
			fields: fields{
				delegate: &mockHttpClient{returns: errors.New("post fail")},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url:  "localhost",
				body: data{A: "1", B: "2"},
			},
			assertion: func(fields fields, args args, err error) {
				assert.Equal(t, errors.New("post fail"), err)
			},
		},
		{
			name: "500 error",
			fields: fields{
				delegate: &mockHttpClient{returns: &http.Response{StatusCode: 500, Body: &mockBody{}}},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url:  "localhost",
				body: data{A: "1", B: "2"},
			},
			assertion: func(fields fields, args args, err error) {
				assert.Equal(t, errors.New("http status code: 500"), err)
			},
		},
		{
			name: "post success",
			fields: fields{
				delegate: &mockHttpClient{returns: &http.Response{StatusCode: 200, Body: &mockBody{}}},
				sdk:      model.Sdk{},
				clock:    clock.Fixed(42),
			},
			args: args{
				url:  "localhost",
				body: data{A: "1", B: "2"},
			},
			assertion: func(fields fields, args args, err error) {
				assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &httpClient{
				delegate: tt.fields.delegate,
				sdk:      tt.fields.sdk,
				clock:    tt.fields.clock,
			}
			err := c.PostObj(tt.args.url, tt.args.body)
			tt.assertion(tt.fields, tt.args, err)
		})
	}
}

type data struct {
	A string `json:"a"`
	B string `json:"b"`
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

type mockBody struct {
	bytes []byte
}

func (m *mockBody) Read(p []byte) (n int, err error) {
	for i, b := range m.bytes {
		p[i] = b
	}
	return len(m.bytes), nil
}

func (m *mockBody) Close() error {
	return nil
}
