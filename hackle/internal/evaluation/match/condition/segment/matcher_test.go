package segment

import (
	"errors"
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConditionMatcher(t *testing.T) {
	matcher := NewConditionMatcher(&mockConditionMatcher{})
	assert.IsType(t, &segmentMatcher{}, matcher.segmentMatcher)
}

func TestConditionMatcher_Matches(t *testing.T) {
	t.Run("when not segment key type then return error", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "age"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"42"}},
		}

		sut := &ConditionMatcher{&mockSegmentMatcher{returns: true}}

		// when
		matches, err := sut.Matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("unsupported target key type [USER_PROPERTY]"), err)
	})

	t.Run("when match value is empty then return false", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeSegment, Name: "SEGMENT"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{}},
		}

		sut := &ConditionMatcher{&mockSegmentMatcher{returns: true}}

		// when
		matches, err := sut.Matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, nil, err)
	})

	t.Run("when segment key is not string type then return error", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeSegment, Name: "SEGMENT"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{42}},
		}

		sut := &ConditionMatcher{&mockSegmentMatcher{returns: true}}

		// when
		matches, err := sut.Matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Contains(t, err.Error(), "segment key")
	})

	t.Run("when segment not found then return error", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{
			W: mocks.CreateWorkspace(),
		}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeSegment, Name: "SEGMENT"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"seg_key"}},
		}

		sut := &ConditionMatcher{&mockSegmentMatcher{returns: true}}

		// when
		matches, err := sut.Matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("segment [seg_key]"), err)
	})

	t.Run("when segment matched then return true", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{
			W: mocks.CreateWorkspace().Segment(model.Segment{Key: "seg_key"}),
		}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeSegment, Name: "SEGMENT"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"seg_key"}},
		}

		sut := &ConditionMatcher{&mockSegmentMatcher{returns: true}}

		// when
		matches, err := sut.Matches(request, context, condition)

		// then
		assert.Equal(t, true, matches)
		assert.Equal(t, nil, err)
	})

	t.Run("NOT_MATCH type", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{
			W: mocks.CreateWorkspace().Segment(model.Segment{Key: "seg_key"}),
		}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeSegment, Name: "SEGMENT"},
			Match: model.TargetMatch{Type: model.MatchTypeNotMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"seg_key"}},
		}

		sut := &ConditionMatcher{&mockSegmentMatcher{returns: true}}

		// when
		matches, err := sut.Matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, nil, err)
	})
}

func Test_segmentMatcher_matches(t *testing.T) {

	segment := func(matcher *mockConditionMatcher, targetConditions ...[]interface{}) model.Segment {
		targets := make([]model.Target, len(targetConditions))
		for i, targetMatches := range targetConditions {
			conditions := make([]model.TargetCondition, len(targetMatches))
			for j, conditionMatch := range targetMatches {
				condition := model.TargetCondition{Key: model.TargetKey{Name: fmt.Sprintf("%d,%d", i, j)}}
				conditions[j] = condition
				matcher.Return(condition, conditionMatch)
			}
			target := model.Target{Conditions: conditions}
			targets[i] = target
		}
		return model.Segment{ID: 42, Key: "seg", Type: model.SegmentTypeUserProperty, Targets: targets}
	}

	t.Run("when target is empty then return false", func(t *testing.T) {
		// given
		matcher := &mockConditionMatcher{returns: make(map[string]interface{})}
		sut := &segmentMatcher{matcher}

		seg := segment(matcher)

		// when
		matches, err := sut.matches(evaluator.SimpleRequest{}, evaluator.NewContext(), seg)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 0, matcher.count)
	})

	t.Run("when any if target matched then return true", func(t *testing.T) {
		// given
		matcher := &mockConditionMatcher{returns: make(map[string]interface{})}
		sut := &segmentMatcher{matcher}

		seg := segment(matcher,
			[]interface{}{true, true, true, false}, // false
			[]interface{}{false},                   // false
			[]interface{}{true, true},              // true
		)

		// when
		matches, err := sut.matches(evaluator.SimpleRequest{}, evaluator.NewContext(), seg)

		// then
		assert.Equal(t, true, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 7, matcher.count)
	})

	t.Run("when all target do not matched then return false", func(t *testing.T) {
		// given
		matcher := &mockConditionMatcher{returns: make(map[string]interface{})}
		sut := &segmentMatcher{matcher}

		seg := segment(matcher,
			[]interface{}{true, true, true, false},
			[]interface{}{false},
			[]interface{}{false, true},
		)

		// when
		matches, err := sut.matches(evaluator.SimpleRequest{}, evaluator.NewContext(), seg)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 6, matcher.count)
	})

	t.Run("when error on matches then return error", func(t *testing.T) {
		// given
		matcher := &mockConditionMatcher{returns: make(map[string]interface{})}
		sut := &segmentMatcher{matcher}

		seg := segment(matcher,
			[]interface{}{true, true, true, false},
			[]interface{}{true, errors.New("match error")},
			[]interface{}{false, true},
		)

		// when
		matches, err := sut.matches(evaluator.SimpleRequest{}, evaluator.NewContext(), seg)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("match error"), err)
		assert.Equal(t, 6, matcher.count)
	})
}

type mockSegmentMatcher struct {
	returns interface{}
}

func (m *mockSegmentMatcher) matches(request evaluator.Request, context evaluator.Context, segment model.Segment) (bool, error) {
	switch r := m.returns.(type) {
	case bool:
		return r, nil
	case error:
		return false, r
	default:
		return false, nil
	}
}

type mockConditionMatcher struct {
	returns map[string]interface{}
	count   int
}

func (m *mockConditionMatcher) Matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	m.count++
	r, ok := m.returns[condition.Key.Name]
	if !ok {
		return false, nil
	}

	switch rr := r.(type) {
	case bool:
		return rr, nil
	case error:
		return false, rr
	default:
		return false, nil
	}
}

func (m *mockConditionMatcher) Return(condition model.TargetCondition, value interface{}) {
	m.returns[condition.Key.Name] = value
}
