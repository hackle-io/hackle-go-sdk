package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/value"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type ConditionMatcher struct {
	abTestMatcher      experimentMatcher
	featureFlagMatcher experimentMatcher
}

func NewConditionMatcher(evaluator evaluator.Evaluator, matcher value.OperatorMatcher) *ConditionMatcher {
	return &ConditionMatcher{
		abTestMatcher:      newAbTestEvaluatorMatcher(evaluator, matcher),
		featureFlagMatcher: newFeatureFlagEvaluatorMatcher(evaluator, matcher),
	}
}

func (m *ConditionMatcher) Matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	switch condition.Key.Type {
	case model.TargetKeyTypeAbTest:
		return m.abTestMatcher.matches(request, context, condition)
	case model.TargetKeyTypeFeatureFlag:
		return m.featureFlagMatcher.matches(request, context, condition)
	}
	return false, fmt.Errorf("unsupported target key type [%s]", condition.Key.Type)
}
