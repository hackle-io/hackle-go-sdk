package delegating

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/contextual"
)

type Evaluator struct {
	evaluators []contextual.Evaluator
}

func NewEvaluator() *Evaluator {
	return &Evaluator{evaluators: make([]contextual.Evaluator, 0)}
}

func (e *Evaluator) Add(evaluator contextual.Evaluator) {
	e.evaluators = append(e.evaluators, evaluator)
}

func (e *Evaluator) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	for _, ce := range e.evaluators {
		if ce.Supports(request) {
			return ce.EvaluateContextually(request, context)
		}
	}
	return nil, fmt.Errorf("unsupported evaluator.Request: %s", request)
}
