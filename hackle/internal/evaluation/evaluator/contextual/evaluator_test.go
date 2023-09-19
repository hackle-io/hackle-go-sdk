package contextual

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestBaseEvaluator_EvaluateContextually(t *testing.T) {

	t.Run("when context contains request then return error - circular evaluation", func(t *testing.T) {
		// given
		sut := NewBase()

		context := evaluator.NewContext()
		req := evaluator.SimpleRequest{K: evaluator.Key{Type: evaluator.TypeExperiment, ID: 1}}
		context.AddRequest(req)

		// when
		_, err := sut.EvaluateContextually(req, context)

		// then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "circular evaluation has occurred")
	})

	t.Run("evaluate", func(t *testing.T) {
		// given
		sut := NewBase()

		context := evaluator.NewContext()
		req := evaluator.SimpleRequest{K: evaluator.Key{Type: evaluator.TypeExperiment, ID: 1}}
		e := &mockEvaluator{}
		sut.Evaluator = e
		eval := evaluator.SimpleEvaluation{R: "contextually"}
		e.On("Evaluate", mock.Anything, mock.Anything).Return(eval, nil)

		// when
		actual, err := sut.EvaluateContextually(req, context)

		// then
		assert.Nil(t, err)
		assert.Equal(t, eval, actual)
		assert.Empty(t, context.Requests())
	})
}

type mockEvaluator struct{ mock.Mock }

func (m *mockEvaluator) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	arguments := m.Called(request, context)
	return arguments.Get(0).(evaluator.Evaluation), arguments.Error(1)
}

func (m *mockEvaluator) Supports(request evaluator.Request) bool {
	arguments := m.Called(request)
	return arguments.Bool(0)
}

func (m *mockEvaluator) EvaluateContextually(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	arguments := m.Called(request, context)
	return arguments.Get(0).(evaluator.Evaluation), arguments.Error(1)
}
