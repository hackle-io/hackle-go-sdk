package evaluator

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContext_Requests(t *testing.T) {

	context := NewContext()
	assert.Len(t, context.Requests(), 0)

	request1 := mockRequest{key: Key{TypeExperiment, 1}}
	assert.False(t, context.Contains(request1))

	context.AddRequest(request1)
	assert.True(t, context.Contains(request1))

	requests1 := context.Requests()
	assert.Len(t, requests1, 1)

	request2 := mockRequest{key: Key{TypeExperiment, 2}}
	assert.False(t, context.Contains(request2))

	context.AddRequest(request2)
	assert.True(t, context.Contains(request2))

	requests2 := context.Requests()
	assert.Len(t, requests2, 2)

	context.RemoveRequest(request2)
	assert.Len(t, context.Requests(), 1)

	context.RemoveRequest(request1)
	assert.Len(t, context.Requests(), 0)

	assert.Len(t, requests1, 1)
	assert.Len(t, requests2, 2)
	assert.False(t, context.Contains(request1))
	assert.False(t, context.Contains(request2))
}

func TestContext_Evaluations(t *testing.T) {
	context := NewContext()
	assert.Len(t, context.Evaluations(), 0)

	context.AddEvaluation(mockEvaluation{})
	assert.Len(t, context.Evaluations(), 1)

	context.AddEvaluation(mockEvaluation{})
	assert.Len(t, context.Evaluations(), 2)
}

type mockEvaluation struct {
	reason      string
	evaluations []Evaluation
}

func (e mockEvaluation) Reason() string {
	return e.reason
}

func (e mockEvaluation) TargetEvaluations() []Evaluation {
	return e.evaluations
}

type mockRequest struct {
	key       Key
	workspace workspace.Workspace
	user      user.HackleUser
}

func (r mockRequest) String() string {
	return "MockRequest"
}

func (r mockRequest) Key() Key {
	return r.key
}

func (r mockRequest) Workspace() workspace.Workspace {
	return r.workspace
}

func (r mockRequest) User() user.HackleUser {
	return r.user
}
