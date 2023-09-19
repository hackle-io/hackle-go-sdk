package contextual

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
)

type Evaluator interface {
	evaluator.Evaluator
	Supports(request evaluator.Request) bool
	EvaluateContextually(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error)
}

func NewBase() *BaseEvaluator {
	return &BaseEvaluator{}
}

type BaseEvaluator struct {
	Evaluator
}

func (e *BaseEvaluator) EvaluateContextually(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	if context.Contains(request) {
		return nil, fmt.Errorf("circular evaluation has occurred %s", append(context.Requests(), request))
	}

	context.AddRequest(request)
	defer context.RemoveRequest(request)

	return e.Evaluate(request, context)
}
