package hackle

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvent(t *testing.T) {

	t.Run("build", func(t *testing.T) {
		event := NewEventBuilder("purchase").
			Value(42.0).
			Property("k1", "v1").
			Property("k2", 2).
			Properties(map[string]interface{}{"k3": true}).
			Properties(nil).
			Build()
		assert.Equal(t, Event{
			key:   "purchase",
			value: 42.0,
			properties: map[string]interface{}{
				"k1": "v1",
				"k2": 2,
				"k3": true,
			},
		}, event)

		assert.Equal(t, "purchase", event.Key())
		assert.Equal(t, 42.0, event.Value())
		assert.Equal(t, map[string]interface{}{
			"k1": "v1",
			"k2": 2,
			"k3": true,
		}, event.Properties())
	})

	t.Run("new", func(t *testing.T) {
		event := NewEvent("purchase")
		assert.Equal(t, Event{
			key:        "purchase",
			value:      0.0,
			properties: map[string]interface{}{},
		}, event)
	})
}
