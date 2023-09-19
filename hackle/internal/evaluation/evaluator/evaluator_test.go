package evaluator

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluator(t *testing.T) {

	request := SimpleRequest{
		K: Key{Type: "TEST", ID: 42},
		W: mocks.CreateWorkspace(),
		U: user.HackleUser{Identifiers: map[string]string{"a": "b"}},
	}

	assert.Equal(t, Key{Type: "TEST", ID: 42}, request.Key())
	assert.NotNil(t, request.Workspace())
	assert.Equal(t, user.HackleUser{Identifiers: map[string]string{"a": "b"}}, request.User())

	evaluation := SimpleEvaluation{
		R: "42",
		E: []Evaluation{},
	}
	assert.Equal(t, "42", evaluation.Reason())
	assert.Equal(t, 0, len(evaluation.TargetEvaluations()))

	assert.Equal(t, "EXPERIMENT", TypeExperiment.String())
	assert.Equal(t, "REMOTE_CONFIG", TypeRemoteConfig.String())
}
