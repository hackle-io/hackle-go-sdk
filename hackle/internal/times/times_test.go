package times

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMillis(t *testing.T) {
	millis := Millis(float64(42 * time.Millisecond))
	assert.Equal(t, 42.0, millis)
}
