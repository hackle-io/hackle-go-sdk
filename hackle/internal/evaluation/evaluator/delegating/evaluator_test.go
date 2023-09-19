package delegating

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluator(t *testing.T) {

	sut := NewEvaluator()

	_, err := sut.Evaluate(evaluator.SimpleRequest{}, evaluator.NewContext())
	assert.NotNil(t, err)

	r1 := evaluator.SimpleRequest{K: evaluator.Key{Type: evaluator.TypeExperiment, ID: 1}}
	e1 := evaluator.SimpleEvaluation{R: "1"}
	evaluator1 := &mockContextualEvaluator{r1, e1}
	sut.Add(evaluator1)
	actual, _ := sut.Evaluate(r1, evaluator.NewContext())
	assert.Equal(t, e1, actual)

	r2 := evaluator.SimpleRequest{K: evaluator.Key{Type: evaluator.TypeExperiment, ID: 2}}
	_, err = sut.Evaluate(r2, evaluator.NewContext())
	assert.NotNil(t, err)
}

type mockContextualEvaluator struct {
	r evaluator.Request
	e evaluator.Evaluation
}

func (m *mockContextualEvaluator) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	return m.e, nil
}

func (m *mockContextualEvaluator) Supports(request evaluator.Request) bool {
	return request.Key() == m.r.Key()
}

func (m *mockContextualEvaluator) EvaluateContextually(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	return m.e, nil
}
