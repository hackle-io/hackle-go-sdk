package experiment

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTargetDeterminer_IsUserInExperimentTarget(t *testing.T) {

	t.Run("when audiences is empty then return true", func(t *testing.T) {
		// given
		sut := targetDeterminer{&mockTargetMatcher{}}
		request := Request{Experiment: model.Experiment{}}

		// when
		actual, err := sut.IsUserInExperimentTarget(request, evaluator.NewContext())

		// then
		assert.Nil(t, err)
		assert.True(t, actual)
	})

	t.Run("when any of audience match then return true", func(t *testing.T) {
		// given
		matcher := &mockTargetMatcher{returns: []interface{}{false, false, false, true, false}}
		request := Request{Experiment: model.Experiment{TargetAudiences: []model.Target{{}, {}, {}, {}, {}}}}

		sut := targetDeterminer{matcher}

		// when
		actual, err := sut.IsUserInExperimentTarget(request, evaluator.NewContext())

		// then
		assert.Nil(t, err)
		assert.True(t, actual)
		assert.Equal(t, 4, matcher.count)
	})

	t.Run("when all audiences do not match then return false", func(t *testing.T) {
		// given
		matcher := &mockTargetMatcher{returns: []interface{}{false, false, false, false, false}}
		request := Request{Experiment: model.Experiment{TargetAudiences: []model.Target{{}, {}, {}, {}, {}}}}

		sut := targetDeterminer{matcher}

		// when
		actual, err := sut.IsUserInExperimentTarget(request, evaluator.NewContext())

		// then
		assert.Nil(t, err)
		assert.False(t, actual)
		assert.Equal(t, 5, matcher.count)
	})

	t.Run("when error on target matches then return error", func(t *testing.T) {
		// given
		matcher := &mockTargetMatcher{returns: []interface{}{false, errors.New("fail"), false, true, false}}
		request := Request{Experiment: model.Experiment{TargetAudiences: []model.Target{{}, {}, {}, {}, {}}}}

		sut := targetDeterminer{matcher}

		// when
		_, err := sut.IsUserInExperimentTarget(request, evaluator.NewContext())

		// then
		assert.NotNil(t, err)
		assert.Equal(t, 2, matcher.count)
	})
}

func TestTargetRuleDeterminer_Determine(t *testing.T) {
	t.Run("when target rule is empty then return false", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{TargetRules: make([]model.TargetRule, 0)}}
		sut := targetRuleDeterminer{}

		// when
		tr, ok, err := sut.Determine(request, evaluator.NewContext())

		// then
		assert.Equal(t, model.TargetRule{}, tr)
		assert.Equal(t, false, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when target rule matched first then return that target rule", func(t *testing.T) {
		// given
		matcher := &mockTargetMatcher{returns: []interface{}{false, false, false, true, false}}
		request := Request{Experiment: model.Experiment{TargetRules: []model.TargetRule{
			{Action: model.Action{BucketID: ref.Int64(1)}},
			{Action: model.Action{BucketID: ref.Int64(2)}},
			{Action: model.Action{BucketID: ref.Int64(3)}},
			{Action: model.Action{BucketID: ref.Int64(4)}},
			{Action: model.Action{BucketID: ref.Int64(5)}},
		}}}
		sut := targetRuleDeterminer{matcher}

		// when
		tr, ok, err := sut.Determine(request, evaluator.NewContext())

		// then
		assert.Equal(t, ref.Int64(4), tr.Action.BucketID)
		assert.True(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, 4, matcher.count)
	})

	t.Run("when all target rule do not match then return false", func(t *testing.T) {
		// given
		matcher := &mockTargetMatcher{returns: []interface{}{false, false, false, false, false}}
		request := Request{Experiment: model.Experiment{TargetRules: []model.TargetRule{
			{Action: model.Action{BucketID: ref.Int64(1)}},
			{Action: model.Action{BucketID: ref.Int64(2)}},
			{Action: model.Action{BucketID: ref.Int64(3)}},
			{Action: model.Action{BucketID: ref.Int64(4)}},
			{Action: model.Action{BucketID: ref.Int64(5)}},
		}}}
		sut := targetRuleDeterminer{matcher}

		// when
		_, ok, err := sut.Determine(request, evaluator.NewContext())

		// then
		assert.False(t, ok)
		assert.Nil(t, err)
		assert.Equal(t, 5, matcher.count)
	})

	t.Run("when error on target matches then return error", func(t *testing.T) {
		// given
		matcher := &mockTargetMatcher{returns: []interface{}{false, false, errors.New("fail"), false, false}}
		request := Request{Experiment: model.Experiment{TargetRules: []model.TargetRule{
			{Action: model.Action{BucketID: ref.Int64(1)}},
			{Action: model.Action{BucketID: ref.Int64(2)}},
			{Action: model.Action{BucketID: ref.Int64(3)}},
			{Action: model.Action{BucketID: ref.Int64(4)}},
			{Action: model.Action{BucketID: ref.Int64(5)}},
		}}}
		sut := targetRuleDeterminer{matcher}

		// when
		_, ok, err := sut.Determine(request, evaluator.NewContext())

		// then
		assert.False(t, ok)
		assert.NotNil(t, err)
		assert.Equal(t, 3, matcher.count)
	})
}

type mockTargetMatcher struct {
	returns []interface{}
	count   int
}

func (m *mockTargetMatcher) Matches(request evaluator.Request, context evaluator.Context, target model.Target) (bool, error) {
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
