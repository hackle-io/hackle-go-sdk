package target

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcher(t *testing.T) {

	t.Run("when condition is empty then return true", func(t *testing.T) {
		factory := &mockConditionMatcherFactory{
			&mockConditionMatcher{returns: []interface{}{}},
		}
		sut := NewMatcher(factory)

		matches, err := sut.Matches(evaluator.SimpleRequest{}, evaluator.NewContext(), model.Target{})
		assert.Equal(t, true, matches)
		assert.Equal(t, nil, err)
	})

	t.Run("when any of condition not matched then return false", func(t *testing.T) {
		matcher := &mockConditionMatcher{returns: []interface{}{true, true, true, false, true}}
		factory := &mockConditionMatcherFactory{
			matcher,
		}
		sut := NewMatcher(factory)

		target := model.Target{
			Conditions: []model.TargetCondition{{}, {}, {}, {}, {}},
		}
		matches, err := sut.Matches(evaluator.SimpleRequest{}, evaluator.NewContext(), target)
		assert.Equal(t, false, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 4, matcher.count)
	})

	t.Run("when all condition matched then return true", func(t *testing.T) {
		matcher := &mockConditionMatcher{returns: []interface{}{true, true, true, true, true}}
		factory := &mockConditionMatcherFactory{
			matcher,
		}
		sut := NewMatcher(factory)

		target := model.Target{
			Conditions: []model.TargetCondition{{}, {}, {}, {}, {}},
		}
		matches, err := sut.Matches(evaluator.SimpleRequest{}, evaluator.NewContext(), target)
		assert.Equal(t, true, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 5, matcher.count)
	})

	t.Run("when error on condition match then return error", func(t *testing.T) {
		matcher := &mockConditionMatcher{returns: []interface{}{true, true, errors.New("condition error"), true, true}}
		factory := &mockConditionMatcherFactory{
			matcher,
		}
		sut := NewMatcher(factory)

		target := model.Target{
			Conditions: []model.TargetCondition{{}, {}, {}, {}, {}},
		}
		matches, err := sut.Matches(evaluator.SimpleRequest{}, evaluator.NewContext(), target)
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("condition error"), err)
		assert.Equal(t, 3, matcher.count)
	})

	t.Run("when cannot get condition matcher then return error", func(t *testing.T) {
		factory := &mockConditionMatcherFactory{errors.New("condition matcher")}
		sut := NewMatcher(factory)

		target := model.Target{
			Conditions: []model.TargetCondition{{}, {}, {}, {}, {}},
		}
		matches, err := sut.Matches(evaluator.SimpleRequest{}, evaluator.NewContext(), target)
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("condition matcher"), err)
	})
}

type mockConditionMatcher struct {
	returns []interface{}
	count   int
}

func (m *mockConditionMatcher) Matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	r := m.returns[m.count]
	m.count++

	switch rr := r.(type) {
	case bool:
		return rr, nil
	case error:
		return false, rr
	default:
		return false, nil
	}
}

type mockConditionMatcherFactory struct {
	matcher interface{}
}

func (m *mockConditionMatcherFactory) Get(keyType model.TargetKeyType) (condition.Matcher, error) {
	switch r := m.matcher.(type) {
	case error:
		return nil, r
	case condition.Matcher:
		return r, nil
	}
	return nil, nil
}
