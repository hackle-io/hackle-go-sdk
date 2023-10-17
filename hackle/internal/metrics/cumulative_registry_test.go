package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistry(t *testing.T) {

	t.Run("counter", func(t *testing.T) {
		counter := NewCumulativeRegistry().Counter("counter", Tags{})
		assert.IsType(t, &CumulativeCounter{}, counter)
	})

	t.Run("timer", func(t *testing.T) {
		timer := NewCumulativeRegistry().Timer("timer", Tags{})
		assert.IsType(t, &CumulativeTimer{}, timer)
	})

	t.Run("close do noting", func(t *testing.T) {
		NewCumulativeRegistry().Close()
	})
}
