package remoteconfig

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"
)

type Evaluation struct {
	reason            string
	targetEvaluations []evaluator.Evaluation
	Parameter         model.RemoteConfigParameter
	ValueID           *int64
	Value             interface{}
	Properties        map[string]interface{}
}

func NewEvaluation(
	request Request,
	context evaluator.Context,
	valueID *int64,
	value interface{},
	reason string,
	builder *properties.Builder,
) Evaluation {
	builder.Add("returnValue", value)
	return Evaluation{
		reason:            reason,
		targetEvaluations: context.Evaluations(),
		Parameter:         request.Parameter,
		ValueID:           valueID,
		Value:             value,
		Properties:        builder.Build(),
	}
}

func NewEvaluationDefault(request Request, context evaluator.Context, reason string, builder *properties.Builder) Evaluation {
	return NewEvaluation(request, context, nil, request.defaultValue, reason, builder)
}

func NewEvaluationOf(
	reason string,
	targetEvaluations []evaluator.Evaluation,
	parameter model.RemoteConfigParameter,
	valueID *int64,
	value interface{},
	properties map[string]interface{},
) Evaluation {
	return Evaluation{
		reason:            reason,
		targetEvaluations: targetEvaluations,
		Parameter:         parameter,
		ValueID:           valueID,
		Value:             value,
		Properties:        properties,
	}
}

func (e Evaluation) Reason() string {
	return e.reason
}

func (e Evaluation) TargetEvaluations() []evaluator.Evaluation {
	return e.targetEvaluations
}
