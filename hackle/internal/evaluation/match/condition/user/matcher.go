package user

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/value"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
)

type ConditionMatcher struct {
	valueResolver ValueResolver
	matcher       value.OperatorMatcher
}

func NewConditionMatcher(matcher value.OperatorMatcher) *ConditionMatcher {
	return &ConditionMatcher{
		valueResolver: &valueResolver{},
		matcher:       matcher,
	}
}

func (m *ConditionMatcher) Matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	userValue, ok, err := m.valueResolver.resolve(request.User(), condition.Key)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return m.matcher.Matches(userValue, condition.Match), nil
}

type ValueResolver interface {
	resolve(user user.HackleUser, key model.TargetKey) (interface{}, bool, error)
}

type valueResolver struct{}

func (r *valueResolver) resolve(user user.HackleUser, key model.TargetKey) (interface{}, bool, error) {
	switch key.Type {
	case model.TargetKeyTypeUserId:
		v, ok := user.Identifiers[key.Name]
		return v, ok, nil
	case model.TargetKeyTypeUserProperty:
		v, ok := user.Properties[key.Name]
		return v, ok, nil
	case model.TargetKeyTypeHackleProperty:
		return nil, false, nil
	}
	return nil, false, fmt.Errorf("unsupported target key type [%s]", key.Type)
}
