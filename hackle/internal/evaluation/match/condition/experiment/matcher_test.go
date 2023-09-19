package experiment

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConditionMatcher(t *testing.T) {
	matcher := NewConditionMatcher(&mockEvaluator{returns: nil}, &mockValueOperatorMatcher{})
	assert.IsType(t, &abTestEvaluatorMatcher{}, matcher.abTestMatcher)
	assert.IsType(t, &featureFlagEvaluatorMatcher{}, matcher.featureFlagMatcher)
}

func TestConditionMatcher_Matches(t *testing.T) {

	ab := &mockExperimentMatcher{}
	ff := &mockExperimentMatcher{}
	sut := &ConditionMatcher{
		abTestMatcher:      ab,
		featureFlagMatcher: ff,
	}

	_, _ = sut.Matches(experiment.Request{}, evaluator.NewContext(), model.TargetCondition{Key: model.TargetKey{Type: model.TargetKeyTypeAbTest}})
	assert.Equal(t, 1, ab.count)
	assert.Equal(t, 0, ff.count)

	_, _ = sut.Matches(experiment.Request{}, evaluator.NewContext(), model.TargetCondition{Key: model.TargetKey{Type: model.TargetKeyTypeFeatureFlag}})
	assert.Equal(t, 1, ab.count)
	assert.Equal(t, 1, ff.count)

	_, err := sut.Matches(experiment.Request{}, evaluator.NewContext(), model.TargetCondition{Key: model.TargetKey{Type: "unsupported"}})
	assert.Equal(t, errors.New("unsupported target key type [unsupported]"), err)
}

type mockExperimentMatcher struct {
	returns bool
	count   int
}

func (m *mockExperimentMatcher) matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	m.count++
	return m.returns, nil
}
