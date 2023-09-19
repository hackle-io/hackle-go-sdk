package experiment

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRequest(t *testing.T) {

	request := NewRequest(
		mocks.CreateWorkspace(),
		user.HackleUser{Identifiers: map[string]string{"a": "B"}},
		model.Experiment{ID: 42, Key: 320, Type: model.ExperimentTypeAbTest},
		"B",
	)

	assert.Equal(t, evaluator.Key{Type: evaluator.TypeExperiment, ID: 42}, request.Key())
	assert.Equal(t, user.HackleUser{Identifiers: map[string]string{"a": "B"}}, request.User())
	assert.Equal(t, "evaluator.Request(type=AB_TEST, key=320)", request.String())

	request2 := NewRequestFrom(request, model.Experiment{ID: 43, Key: 321, Type: model.ExperimentTypeAbTest})
	assert.Equal(t, evaluator.Key{Type: evaluator.TypeExperiment, ID: 43}, request2.Key())
}
