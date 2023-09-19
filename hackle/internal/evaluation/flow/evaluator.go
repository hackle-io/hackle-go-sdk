package flow

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
)

type Evaluator interface {
	Evaluate(request evaluator.Request, context evaluator.Context, nextFlow EvaluationFlow) (evaluator.Evaluation, bool, error)
}
