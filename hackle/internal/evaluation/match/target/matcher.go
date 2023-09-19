package target

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type Matcher interface {
	Matches(request evaluator.Request, context evaluator.Context, target model.Target) (bool, error)
}

func NewMatcher(conditionMatcherFactory condition.MatcherFactory) Matcher {
	return &matcher{
		conditionMatcherFactory: conditionMatcherFactory,
	}
}

type matcher struct {
	conditionMatcherFactory condition.MatcherFactory
}

func (m *matcher) Matches(request evaluator.Request, context evaluator.Context, target model.Target) (bool, error) {
	for _, it := range target.Conditions {
		matches, err := m.matches(request, context, it)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, nil
		}
	}
	return true, nil
}

func (m *matcher) matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	conditionMatcher, err := m.conditionMatcherFactory.Get(condition.Key.Type)
	if err != nil {
		return false, err
	}
	return conditionMatcher.Matches(request, context, condition)
}
