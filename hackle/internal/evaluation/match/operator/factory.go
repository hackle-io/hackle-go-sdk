package operator

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type MatcherFactory interface {
	Get(operator model.TargetOperator) (Matcher, bool)
}

func NewMatcherFactory() MatcherFactory {
	return &matcherFactory{
		matchers: map[model.TargetOperator]Matcher{
			model.OperatorIn:         &InMatcher{},
			model.OperatorContains:   &containsMatcher{},
			model.OperatorStartsWith: &startsWithMatcher{},
			model.OperatorEndsWith:   &endsWithMatcher{},
			model.OperatorGT:         &greaterThanMatcher{},
			model.OperatorGTE:        &greaterThanOrEqualToMatcher{},
			model.OperatorLT:         &lessThanMatcher{},
			model.OperatorLTE:        &lessThanOrEqualToMatcher{},
		},
	}
}

type matcherFactory struct {
	matchers map[model.TargetOperator]Matcher
}

func (f *matcherFactory) Get(operator model.TargetOperator) (Matcher, bool) {
	matcher, ok := f.matchers[operator]
	return matcher, ok
}
