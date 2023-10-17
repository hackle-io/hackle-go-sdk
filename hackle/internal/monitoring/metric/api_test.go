package metric

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecordAPI(t *testing.T) {

	metrics.AddRegistry(metrics.NewCumulativeRegistry())

	clock := &mockClock{returns: []int64{100, 142, 1000, 1420}}

	RecordAPI(OperationGetWorkspace, metrics.NewTimerSample(clock), true)
	RecordAPI(OperationPostEvents, metrics.NewTimerSample(clock), false)

	workspaceTimer := metrics.NewTimer("api.call", metrics.Tags{"operation": "get.workspace", "success": "true"})
	assert.Equal(t, int64(1), workspaceTimer.Count())
	assert.Equal(t, int64(42), workspaceTimer.Sum())

	eventTimer := metrics.NewTimer("api.call", metrics.Tags{"operation": "post.events", "success": "false"})
	assert.Equal(t, int64(1), eventTimer.Count())
	assert.Equal(t, int64(420), eventTimer.Sum())
}

type mockClock struct {
	returns []int64
	count   int
}

func (m *mockClock) CurrentMillis() int64 {
	return 0
}

func (m *mockClock) Tick() int64 {
	t := m.returns[m.count]
	m.count++
	return t
}
