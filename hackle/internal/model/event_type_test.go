package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUndefinedEvent(t *testing.T) {
	event := NewUndefinedEvent("42")
	assert.Equal(t, EventType{ID: 0, Key: "42"}, event)
}
