package remoteconfig

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {

	request := NewRequest(
		mocks.CreateWorkspace(),
		user.HackleUser{Identifiers: map[string]string{"a": "b"}},
		model.RemoteConfigParameter{ID: 42, Key: "rc"},
		types.String,
		"default",
	)
	assert.Equal(t, evaluator.Key{Type: evaluator.TypeRemoteConfig, ID: 42}, request.Key())
	assert.Equal(t, "evaluator.Request(type=REMOTE_CONFIG, key=rc)", request.String())
}
