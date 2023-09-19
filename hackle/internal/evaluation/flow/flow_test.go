package flow

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluationFlow(t *testing.T) {
	end := &End{}
	evaluation, ok, err := end.Evaluate(evaluator.SimpleRequest{}, evaluator.NewContext())
	assert.Equal(t, nil, evaluation)
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, err)

	decision := &Decision{
		Evaluator: &mockEvaluator{returns: evaluator.SimpleEvaluation{R: "42"}},
		NextFlow:  end,
	}
	evaluation, ok, err = decision.Evaluate(evaluator.SimpleRequest{}, evaluator.NewContext())
	assert.Equal(t, evaluator.SimpleEvaluation{R: "42"}, evaluation)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, err)
}

func TestNewEvaluationFlow(t *testing.T) {

	decisionWith := func(f EvaluationFlow, e Evaluator) EvaluationFlow {
		decision, ok := f.(*Decision)
		assert.True(t, ok)
		assert.Equal(t, e, decision.Evaluator)
		return decision.NextFlow
	}

	end := func(f EvaluationFlow) {
		assert.IsType(t, f, &End{})
	}

	e1 := &mockEvaluator{name: "1"}
	e2 := &mockEvaluator{name: "2"}
	e3 := &mockEvaluator{name: "3"}

	flow := NewEvaluationFlow(e1, e2, e3)

	flow = decisionWith(flow, e1)
	flow = decisionWith(flow, e2)
	flow = decisionWith(flow, e3)
	end(flow)
}

type mockEvaluator struct {
	name    string
	returns evaluator.Evaluation
}

func (m *mockEvaluator) Evaluate(request evaluator.Request, context evaluator.Context, nextFlow EvaluationFlow) (evaluator.Evaluation, bool, error) {
	return m.returns, true, nil
}
