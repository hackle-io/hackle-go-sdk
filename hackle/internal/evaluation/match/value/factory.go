package value

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
)

type MatcherFactory interface {
	Get(valueType types.ValueType) (Matcher, bool)
}

func NewMatcherFactory() MatcherFactory {
	return &matcherFactory{
		matchers: map[types.ValueType]Matcher{
			types.String:  &stringMatcher{},
			types.Number:  &numberMatcher{},
			types.Bool:    &boolMatcher{},
			types.Version: &versionMatcher{},
			types.Json:    &stringMatcher{},
		},
	}
}

type matcherFactory struct {
	matchers map[types.ValueType]Matcher
}

func (f *matcherFactory) Get(valueType types.ValueType) (Matcher, bool) {
	matcher, ok := f.matchers[valueType]
	return matcher, ok
}
