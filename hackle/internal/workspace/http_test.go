package workspace

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewHttpFetcher(t *testing.T) {
	fetcher := NewHttpFetcher("localhost", &mockHttpClient{})
	assert.IsType(t, &httpFetcher{}, fetcher)
	assert.Equal(t, "localhost/api/v2/workspaces", fetcher.(*httpFetcher).url)
}

func TestHttpFetcher_Fetch(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		fetcher := NewHttpFetcher("localhost", &mockHttpClient{returns: WorkspaceDTO{}})
		ws, err := fetcher.Fetch()
		assert.NotNil(t, ws)
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		fetcher := NewHttpFetcher("localhost", &mockHttpClient{err: errors.New("fetch fail")})
		ws, err := fetcher.Fetch()
		assert.Nil(t, ws)
		assert.Equal(t, errors.New("fetch fail"), err)
	})
}

type mockHttpClient struct {
	returns interface{}
	err     error
}

func (m *mockHttpClient) GetObj(url string, out interface{}) error {
	if m.err != nil {
		return m.err
	}
	out = m.returns
	return nil
}

func (m *mockHttpClient) PostObj(url string, body interface{}) error {
	if m.err != nil {
		return m.err
	}
	return nil
}
