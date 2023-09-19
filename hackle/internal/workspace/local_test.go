package workspace

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetch(t *testing.T) {

	fetcher1 := NewFileFetcher("../../../testdata/workspace_config.json")
	defer fetcher1.Close()

	ws, _ := fetcher1.Fetch()
	assert.NotNil(t, ws)

	fetcher2 := NewFileFetcher("invalid")
	_, ok := fetcher2.Fetch()
	assert.Equal(t, false, ok)
}
