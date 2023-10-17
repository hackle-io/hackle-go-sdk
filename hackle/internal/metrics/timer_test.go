package metrics

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimerBuilder_Register(t *testing.T) {

	timer := NewTimerBuilder("timer").
		Tags(Tags{"a": "1", "b": "2"}).
		Tag("c", "3").
		Register(NewCumulativeRegistry())

	assert.Equal(t,
		NewID("timer", Tags{"a": "1", "b": "2", "c": "3"}, TypeTimer),
		timer.ID(),
	)
}

func TestTimerSample(t *testing.T) {
	timer := NewCumulativeRegistry().Timer("timer", Tags{})
	sample := NewTimerSample(clock.System)
	time.Sleep(1000 * time.Millisecond)
	sample.Stop(timer)

	assert.InEpsilon(t, 1000*time.Millisecond, timer.Sum(), 0.01)
}
