package condition

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type MatcherFactory interface {
	Get(keyType model.TargetKeyType) (Matcher, error)
}

func NewMatcherFactory(matchers map[model.TargetKeyType]Matcher) MatcherFactory {
	return &matcherFactory{
		matchers: matchers,
	}
}

type matcherFactory struct {
	matchers map[model.TargetKeyType]Matcher
}

func (f *matcherFactory) Get(keyType model.TargetKeyType) (Matcher, error) {
	matcher, ok := f.matchers[keyType]
	if !ok {
		return nil, fmt.Errorf("unsupported TargetKeyType [%s]", keyType)
	}
	return matcher, nil
}
