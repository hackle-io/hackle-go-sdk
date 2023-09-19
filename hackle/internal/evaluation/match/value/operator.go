package value

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/operator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
)

type OperatorMatcher interface {
	Matches(userValue interface{}, match model.TargetMatch) bool
}

func NewOperatorMatcher() OperatorMatcher {
	return &operatorMatcher{
		valueMatcherFactory:    NewMatcherFactory(),
		operatorMatcherFactory: operator.NewMatcherFactory(),
	}
}

type operatorMatcher struct {
	valueMatcherFactory    MatcherFactory
	operatorMatcherFactory operator.MatcherFactory
}

func (m *operatorMatcher) Matches(userValue interface{}, match model.TargetMatch) bool {
	valueMatcher, ok := m.valueMatcherFactory.Get(match.ValueType)
	if !ok {
		return false
	}

	operatorMatcher, ok := m.operatorMatcherFactory.Get(match.Operator)
	if !ok {
		return false
	}

	matches := m.matches(userValue, match, valueMatcher, operatorMatcher)
	return match.Type.Matches(matches)
}

func (m *operatorMatcher) matches(
	userValue interface{},
	match model.TargetMatch,
	valueMatcher Matcher,
	operatorMatcher operator.Matcher,
) bool {
	if userValues, ok := types.AsArray(userValue); ok {
		return m.arrayMatches(userValues, match, valueMatcher, operatorMatcher)
	} else {
		return m.singleMatches(userValue, match, valueMatcher, operatorMatcher)
	}
}

func (m *operatorMatcher) singleMatches(
	userValue interface{},
	match model.TargetMatch,
	valueMatcher Matcher,
	operatorMatcher operator.Matcher,
) bool {
	for _, matchValue := range match.Values {
		if valueMatcher.Matches(operatorMatcher, userValue, matchValue) {
			return true
		}
	}
	return false
}

func (m *operatorMatcher) arrayMatches(
	userValues []interface{},
	match model.TargetMatch,
	valueMatcher Matcher,
	operatorMatcher operator.Matcher,
) bool {
	for _, userValue := range userValues {
		if m.singleMatches(userValue, match, valueMatcher, operatorMatcher) {
			return true
		}
	}
	return false
}
