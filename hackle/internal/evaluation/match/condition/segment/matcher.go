package segment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type ConditionMatcher struct {
	segmentMatcher SegmentMatcher
}

func NewConditionMatcher(userConditionMatcher condition.Matcher) *ConditionMatcher {
	return &ConditionMatcher{
		segmentMatcher: &segmentMatcher{
			userConditionMatcher: userConditionMatcher,
		},
	}
}

func (m *ConditionMatcher) Matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	if condition.Key.Type != model.TargetKeyTypeSegment {
		return false, fmt.Errorf("unsupported target key type [%s]", condition.Key.Type)
	}
	for _, value := range condition.Match.Values {
		matches, err := m.valueMatches(request, context, value)
		if err != nil {
			return false, err
		}
		if matches {
			return condition.Match.Type.Matches(matches), nil
		}
	}
	return false, nil
}

func (m *ConditionMatcher) valueMatches(request evaluator.Request, context evaluator.Context, value interface{}) (bool, error) {
	segmentKey, ok := value.(string)
	if !ok {
		return false, fmt.Errorf("segment key [%s]", value)
	}
	segment, ok := request.Workspace().GetSegment(segmentKey)
	if !ok {
		return false, fmt.Errorf("segment [%s]", segmentKey)
	}
	return m.segmentMatcher.matches(request, context, segment)
}

//goland:noinspection GoNameStartsWithPackageName
type SegmentMatcher interface {
	matches(request evaluator.Request, context evaluator.Context, segment model.Segment) (bool, error)
}

type segmentMatcher struct {
	userConditionMatcher condition.Matcher
}

func (m *segmentMatcher) matches(request evaluator.Request, context evaluator.Context, segment model.Segment) (bool, error) {
	for _, target := range segment.Targets {
		matches, err := m.targetMatches(request, context, target)
		if err != nil {
			return false, err
		}
		if matches {
			return true, nil
		}
	}
	return false, nil
}

func (m *segmentMatcher) targetMatches(request evaluator.Request, context evaluator.Context, target model.Target) (bool, error) {
	for _, c := range target.Conditions {
		matches, err := m.userConditionMatcher.Matches(request, context, c)
		if err != nil {
			return false, err
		}
		if !matches {
			return false, err
		}
	}
	return true, nil
}
