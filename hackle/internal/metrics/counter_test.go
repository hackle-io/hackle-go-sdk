package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseCounter(t *testing.T) {

	t.Run("new with self", func(t *testing.T) {
		counter := NewBaseCounter(NewCumulativeCounter(NewID("counter", Tags{}, TypeCounter)))
		measurements := counter.Measure()
		assert.Equal(t, 1, len(measurements))
		assert.Equal(t, FieldCount, measurements[0].Field)
	})
}

func TestCounterBuilder(t *testing.T) {

	counter := NewCounterBuilder("test_counter").
		Tags(Tags{"a": "1", "b": "2"}).
		Tag("c", "3").
		Register(NewCumulativeRegistry())

	assert.Equal(t,
		NewID("test_counter", Tags{"a": "1", "b": "2", "c": "3"}, TypeCounter),
		counter.ID())
}
