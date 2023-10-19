package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	nethttp "net/http"
	"sync"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher("localhost", &mockHttpClient{}).(*dispatcher)
	assert.Equal(t, "localhost/api/v2/events", d.url)
}

func Test_dispatcher_Dispatch(t *testing.T) {

	event := NewTrackEvent(
		model.EventType{ID: 42, Key: "my_key"},
		event{key: "my_key"},
		user.NewHackleUserBuilder().Identifier("$id", "id").Build(),
		4200,
	)

	type fields struct {
		url        string
		httpClient *mockHttpClient
		wg         *sync.WaitGroup
	}
	type args struct {
		userEvents []UserEvent
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		assertion func(fields fields)
	}{
		{
			name: "failed to create request",
			fields: fields{
				url:        ":",
				httpClient: &mockHttpClient{},
				wg:         &sync.WaitGroup{},
			},
			args: args{
				userEvents: []UserEvent{event},
			},
			assertion: func(fields fields) {
				fields.wg.Wait()
				assert.Equal(t, false, fields.httpClient.Called())
			},
		},
		{
			name: "failed to http execute",
			fields: fields{
				url: "localhost",
				httpClient: &mockHttpClient{
					err: errors.New("http fail"),
				},
				wg: &sync.WaitGroup{},
			},
			args: args{
				userEvents: []UserEvent{event},
			},
			assertion: func(fields fields) {
				fields.wg.Wait()
				assert.Equal(t, true, fields.httpClient.Called())
			},
		},
		{
			name: "http not success",
			fields: fields{
				url: "localhost",
				httpClient: &mockHttpClient{
					res: &nethttp.Response{
						StatusCode: 500,
						Body:       ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
					},
				},
				wg: &sync.WaitGroup{},
			},
			args: args{
				userEvents: []UserEvent{event},
			},
			assertion: func(fields fields) {
				fields.wg.Wait()
				assert.Equal(t, true, fields.httpClient.Called())
			},
		},
		{
			name: "success",
			fields: fields{
				url: "localhost",
				httpClient: &mockHttpClient{
					res: &nethttp.Response{
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
					},
				},
				wg: &sync.WaitGroup{},
			},
			args: args{
				userEvents: []UserEvent{event},
			},
			assertion: func(fields fields) {
				fields.wg.Wait()
				assert.Equal(t, true, fields.httpClient.Called())
				body, _ := ioutil.ReadAll(fields.httpClient.req.Body)
				var dto PayloadDTO
				_ = json.Unmarshal(body, &dto)
				assert.Equal(t, 1, len(dto.TrackEvents))
				assert.Equal(t, "my_key", dto.TrackEvents[0].EventTypeKey)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &dispatcher{
				url:        tt.fields.url,
				httpClient: tt.fields.httpClient,
				wg:         tt.fields.wg,
			}
			d.Dispatch(tt.args.userEvents)
			tt.assertion(tt.fields)
		})
	}
}

func Test_dispatcher_Close(t *testing.T) {

	t.Run("wait dispatch", func(t *testing.T) {
		httpClient := &mockHttpClient{
			res: &nethttp.Response{
				StatusCode: 500,
				Body:       ioutil.NopCloser(bytes.NewReader(make([]byte, 0))),
			},
			delay: 100 * time.Millisecond,
		}

		sut := NewDispatcher("localhost", httpClient)

		sut.Dispatch(make([]UserEvent, 0))

		assert.Equal(t, false, httpClient.Called())
		sut.Close()
		assert.Equal(t, true, httpClient.Called())
	})
}

type mockHttpClient struct {
	mu    sync.Mutex
	req   *nethttp.Request
	res   *nethttp.Response
	delay time.Duration
	err   error
}

func (m *mockHttpClient) Execute(req *nethttp.Request) (*nethttp.Response, error) {
	time.Sleep(m.delay)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.req = req

	if m.err != nil {
		return nil, m.err
	}
	return m.res, nil
}

func (m *mockHttpClient) Called() bool {
	return m.req != nil
}
