package event

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTrackEvent(t *testing.T) {

	e := NewTrackEvent(
		model.EventType{ID: 42, Key: "my_key"},
		event{key: "my_key"},
		user.HackleUser{Identifiers: map[string]string{"a": "b"}},
		4200,
	)

	assert.NotEmpty(t, e.InsertID())
	assert.Equal(t, int64(4200), e.Timestamp())
	assert.Equal(t, user.HackleUser{Identifiers: map[string]string{"a": "b"}}, e.User())
	assert.Equal(t, event{key: "my_key"}, e.Event)
	assert.Equal(t, model.EventType{ID: 42, Key: "my_key"}, e.EventType)
}
