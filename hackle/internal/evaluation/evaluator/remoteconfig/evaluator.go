package remoteconfig

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/contextual"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
)

type Evaluator interface {
	EvaluateRemoteConfig(request Request, context evaluator.Context) (Evaluation, error)
}

type ContextualEvaluator interface {
	Evaluator
	contextual.Evaluator
}

func NewEvaluator(determiner TargetRuleDeterminer) ContextualEvaluator {
	e := &remoteConfigEvaluator{
		BaseEvaluator: contextual.NewBase(),
		determiner:    determiner,
	}
	e.Evaluator = e
	return e
}

type remoteConfigEvaluator struct {
	*contextual.BaseEvaluator
	determiner TargetRuleDeterminer
}

func (e *remoteConfigEvaluator) Supports(request evaluator.Request) bool {
	_, ok := request.(Request)
	return ok
}

func (e *remoteConfigEvaluator) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	remoteConfigRequest, ok := request.(Request)
	if !ok {
		return nil, fmt.Errorf("unsupported evaluator request: %T (expected: remoteconfig.Request)", request)
	}
	return e.EvaluateRemoteConfig(remoteConfigRequest, context)
}

func (e *remoteConfigEvaluator) EvaluateRemoteConfig(request Request, context evaluator.Context) (Evaluation, error) {
	propertiesBuilder := properties.NewBuilder()
	propertiesBuilder.Add("requestValueType", request.requiredType.String())
	propertiesBuilder.Add("requestDefaultValue", request.defaultValue)

	parameter := request.Parameter
	if _, ok := request.User().Identifiers[parameter.IdentifierType]; !ok {
		return NewEvaluationDefault(request, context, decision.ReasonIdentifierNotFound, propertiesBuilder), nil
	}

	targetRule, ok, err := e.determiner.Determine(request, context)
	if err != nil {
		return Evaluation{}, err
	}
	if ok {
		propertiesBuilder.Add("targetRuleKey", targetRule.Key)
		propertiesBuilder.Add("targetRuleName", targetRule.Name)
		return newEvaluation(request, context, targetRule.Value, decision.ReasonTargetRuleMatch, propertiesBuilder), nil
	}

	return newEvaluation(request, context, parameter.DefaultValue, decision.ReasonDefaultRule, propertiesBuilder), nil
}

func newEvaluation(
	request Request,
	context evaluator.Context,
	parameterValue model.RemoteConfigValue,
	reason string,
	builder *properties.Builder,
) Evaluation {
	if value, ok := rawValue(request, parameterValue); ok {
		return NewEvaluation(request, context, &parameterValue.ID, value, reason, builder)
	}
	return NewEvaluationDefault(request, context, decision.ReasonTypeMismatch, builder)
}

func rawValue(request Request, value model.RemoteConfigValue) (interface{}, bool) {
	switch request.requiredType {
	case types.String:
		if s, ok := value.RawValue.(string); ok {
			return s, true
		}
	case types.Number:
		if types.IsNumber(value.RawValue) {
			return types.AsNumber(value.RawValue)
		}
	case types.Bool:
		if b, ok := value.RawValue.(bool); ok {
			return b, true
		}
	}
	return nil, false
}
