package event

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	d := NewDispatcher("localhost", &mockHttpClient{}).(*dispatcher)
	assert.Equal(t, "localhost/api/v2/events", d.url)
}

func TestDispatcher_Dispatch(t *testing.T) {

	t.Run("dispatch async", func(t *testing.T) {
		httpClient := &mockHttpClient{delay: 100 * time.Millisecond}
		wg := &sync.WaitGroup{}
		sut := &dispatcher{
			url:        "localhost",
			httpClient: httpClient,
			wg:         wg,
		}

		events := make([]UserEvent, 42)
		sut.Dispatch(events)

		assert.Equal(t, 0, len(httpClient.Posts()))
		wg.Wait()
		assert.Equal(t, 1, len(httpClient.Posts()))
	})

	t.Run("dispatch failed", func(t *testing.T) {
		httpClient := &mockHttpClient{delay: 100 * time.Millisecond, err: errors.New("http error")}
		wg := &sync.WaitGroup{}
		sut := &dispatcher{
			url:        "localhost",
			httpClient: httpClient,
			wg:         wg,
		}

		events := make([]UserEvent, 42)
		sut.Dispatch(events)

		assert.Equal(t, 0, len(httpClient.Posts()))
		wg.Wait()
		assert.Equal(t, 0, len(httpClient.Posts()))
	})
}

func TestDispatcher_Close(t *testing.T) {
	t.Run("wait for dispatching", func(t *testing.T) {
		httpClient := &mockHttpClient{delay: 100 * time.Millisecond}
		wg := &sync.WaitGroup{}
		sut := &dispatcher{
			url:        "localhost",
			httpClient: httpClient,
			wg:         wg,
		}

		for i := 0; i < 100; i++ {
			sut.Dispatch(make([]UserEvent, 1))
		}
		assert.Equal(t, 0, len(httpClient.Posts()))
		sut.Close()
		assert.Equal(t, 100, len(httpClient.Posts()))
	})
}

type mockHttpClient struct {
	mu    sync.Mutex
	posts []interface{}
	delay time.Duration
	err   error
}

func (m *mockHttpClient) GetObj(url string, out interface{}) error {
	time.Sleep(m.delay)
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockHttpClient) PostObj(url string, body interface{}) error {
	time.Sleep(m.delay)
	if m.err != nil {
		return m.err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.posts = append(m.posts, body)
	return nil
}

func (m *mockHttpClient) Posts() []interface{} {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.posts
}
