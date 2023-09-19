package flow

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
)

type EvaluationFlow interface {
	Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, bool, error)
}

func NewEvaluationFlow(evaluators ...Evaluator) EvaluationFlow {
	var flow EvaluationFlow = &End{}
	for i := len(evaluators) - 1; i >= 0; i-- {
		flow = &Decision{Evaluator: evaluators[i], NextFlow: flow}
	}
	return flow
}

type End struct{}

func (e *End) Evaluate(evaluator.Request, evaluator.Context) (evaluator.Evaluation, bool, error) {
	return nil, false, nil
}

type Decision struct {
	Evaluator Evaluator
	NextFlow  EvaluationFlow
}

func (d *Decision) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, bool, error) {
	return d.Evaluator.Evaluate(request, context, d.NextFlow)
}
