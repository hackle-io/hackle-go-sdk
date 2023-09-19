package clock

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSystemClock_CurrentMillis(t *testing.T) {
	clock := SystemClock{}
	start := clock.CurrentMillis()
	time.Sleep(500 * time.Millisecond)
	end := clock.CurrentMillis()
	assert.InEpsilon(t, 500, end-start, 0.01)
}

func TestSystemClock_Tick(t *testing.T) {
	clock := SystemClock{}
	start := clock.Tick()
	time.Sleep(500000000 * time.Nanosecond)
	end := clock.Tick()
	assert.InEpsilon(t, 500000000, end-start, 0.01)
}

func TestFixed(t *testing.T) {
	clock := Fixed(42)
	assert.Equal(t, int64(42), clock.CurrentMillis())
	assert.Equal(t, int64(42), clock.Tick())
}
