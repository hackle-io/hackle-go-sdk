package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/contextual"
)

type Evaluator interface {
	EvaluateExperiment(request Request, context evaluator.Context) (Evaluation, error)
}

type ContextualEvaluator interface {
	Evaluator
	contextual.Evaluator
}

func NewEvaluator(factory EvaluationFlowFactory) ContextualEvaluator {
	e := &experimentEvaluator{
		BaseEvaluator: contextual.NewBase(),
		flowFactory:   factory,
	}
	e.Evaluator = e
	return e
}

type experimentEvaluator struct {
	*contextual.BaseEvaluator
	flowFactory EvaluationFlowFactory
}

func (e *experimentEvaluator) Supports(request evaluator.Request) bool {
	_, ok := request.(Request)
	return ok
}

func (e *experimentEvaluator) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	experimentRequest, ok := request.(Request)
	if !ok {
		return nil, fmt.Errorf("unsupported evaluator request: %T (expected: experiment.Request)", request)
	}
	return e.EvaluateExperiment(experimentRequest, context)
}

func (e *experimentEvaluator) EvaluateExperiment(request Request, context evaluator.Context) (Evaluation, error) {
	flow, err := e.flowFactory.Get(request.Experiment.Type)
	if err != nil {
		return Evaluation{}, err
	}
	evaluation, ok, err := flow.Evaluate(request, context)
	if err != nil {
		return Evaluation{}, err
	}
	if !ok {
		return NewEvaluationDefault(request, context, decision.ReasonTrafficNotAllocated)
	}
	experimentEvaluation, ok := evaluation.(Evaluation)
	if !ok {
		return Evaluation{}, fmt.Errorf("unexpected evaluation: %T (expected: experiment.Evaluation)", evaluation)
	}
	return experimentEvaluation, nil
}
