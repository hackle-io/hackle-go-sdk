package experiment

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newAbTestEvaluatorMatcher(t *testing.T) {
	matcher := newAbTestEvaluatorMatcher(&mockEvaluator{}, &mockValueOperatorMatcher{})
	assert.Equal(t, matcher.evaluatorMatcher, matcher)
}

func TestAbTestEvaluatorMatcher(t *testing.T) {

	type fields struct {
		evaluator *mockEvaluator
		matcher   *mockValueOperatorMatcher
	}

	t.Run("when key is not number then return error", func(t *testing.T) {
		// given
		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{},
			&mockValueOperatorMatcher{},
		)

		request := experiment.Request{
			Experiment: model.Experiment{
				ID:   42,
				Type: model.ExperimentTypeAbTest,
			},
		}
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "string"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		// when
		matches, err := sut.matches(request, evaluator.NewContext(), condition)

		// then
		assert.Equal(t, false, matches)
		assert.NotNil(t, err)
	})

	t.Run("when experiment not found then return false", func(t *testing.T) {
		// given
		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{},
			&mockValueOperatorMatcher{},
		)

		exp := model.Experiment{
			ID:   42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace(),
			user.HackleUser{},
			exp,
			"A",
		)
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		// when
		matches, err := sut.matches(request, evaluator.NewContext(), condition)

		// then
		assert.Equal(t, false, matches)
		assert.Nil(t, err)
	})

	t.Run("when error on target evaluate then return error", func(t *testing.T) {
		// given
		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: errors.New("target evaluate error")},
			&mockValueOperatorMatcher{},
		)

		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("target evaluate error"), err)
	})

	t.Run("when target evaluated evaluation is not experiment evaluation type then return error", func(t *testing.T) {
		// given
		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: evaluator.SimpleEvaluation{}},
			&mockValueOperatorMatcher{},
		)

		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Contains(t, err.Error(), "unexpected evaluation")
	})

	t.Run("when not experiment request then use evaluated evaluation directly", func(t *testing.T) {
		// given
		targetEvaluation := experiment.Evaluation{VariationKey: "target evaluation"}
		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: targetEvaluation},
			&mockValueOperatorMatcher{},
		)

		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}

		request := evaluator.SimpleRequest{
			W: mocks.CreateWorkspace().Experiment(exp),
			U: user.HackleUser{},
		}
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, targetEvaluation, context.Evaluations()[0])
		assert.Equal(t, false, matches)
		assert.Nil(t, err)
	})

	t.Run("when target evaluation is not traffic allocated reason then use evaluation directly", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonOverridden)

		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: evaluation},
			&mockValueOperatorMatcher{},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, evaluation, context.Evaluations()[0])
		assert.Equal(t, false, matches)
		assert.Nil(t, err)
	})

	t.Run("when target evaluation is traffic allocated reason then replace reason", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonTrafficAllocated)

		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: evaluation},
			&mockValueOperatorMatcher{},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, decision.ReasonTrafficAllocatedByTargeting, context.Evaluations()[0].Reason())
		assert.Equal(t, false, matches)
		assert.Nil(t, err)
	})

	t.Run("when target evaluation is matched reason then match variation", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonTrafficAllocated)

		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: evaluation},
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, true, matches)
		assert.Nil(t, err)
	})

	t.Run("when target evaluation is not matched reason then return false", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonExperimentDraft)

		sut := newAbTestEvaluatorMatcher(
			&mockEvaluator{returns: evaluation},
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Nil(t, err)
	})

	t.Run("when already evaluated experiment then not evaluate again", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeAbTest,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().Experiment(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"A"}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonOverridden)

		context.AddEvaluation(evaluation)

		targetEvaluator := &mockEvaluator{returns: evaluation}
		sut := newAbTestEvaluatorMatcher(
			targetEvaluator,
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, true, matches)
		assert.Nil(t, err)
		assert.Equal(t, 0, targetEvaluator.count)
	})
}

func Test_newFeatureFlagEvaluatorMatcher(t *testing.T) {
	matcher := newFeatureFlagEvaluatorMatcher(&mockEvaluator{}, &mockValueOperatorMatcher{})
	assert.Equal(t, matcher.evaluatorMatcher, matcher)
}

func TestFeatureFlatEvaluatorMatcher_matches(t *testing.T) {

	t.Run("when key is not number then return error", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeFeatureFlag,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().FeatureFlag(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeFeatureFlag, Name: "string"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.Bool, Values: []interface{}{true}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonDefaultRule)

		targetEvaluator := &mockEvaluator{returns: evaluation}
		sut := newFeatureFlagEvaluatorMatcher(
			targetEvaluator,
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, errors.New("invalid key [FEATURE_FLAG, string]"), err)
	})

	t.Run("when feature flag not found then return false", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeFeatureFlag,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace(),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeFeatureFlag, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.Bool, Values: []interface{}{true}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonDefaultRule)

		targetEvaluator := &mockEvaluator{returns: evaluation}
		sut := newFeatureFlagEvaluatorMatcher(
			targetEvaluator,
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, false, matches)
		assert.Equal(t, nil, err)
	})

	t.Run("when already evaluated then do not evaluate again", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeFeatureFlag,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().FeatureFlag(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeFeatureFlag, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.Bool, Values: []interface{}{true}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonDefaultRule)
		context.AddEvaluation(evaluation)

		targetEvaluator := &mockEvaluator{returns: evaluation}
		sut := newFeatureFlagEvaluatorMatcher(
			targetEvaluator,
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, true, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 0, targetEvaluator.count)
	})

	t.Run("new evaluation", func(t *testing.T) {
		// given
		exp := model.Experiment{
			Key:  42,
			Type: model.ExperimentTypeFeatureFlag,
		}
		request := experiment.NewRequest(
			mocks.CreateWorkspace().FeatureFlag(exp),
			user.HackleUser{},
			exp,
			"A",
		)
		context := evaluator.NewContext()
		condition := model.TargetCondition{
			Key:   model.TargetKey{Type: model.TargetKeyTypeFeatureFlag, Name: "42"},
			Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.Bool, Values: []interface{}{true}},
		}

		evaluation, _ := experiment.NewEvaluationDefault(request, context, decision.ReasonDefaultRule)

		targetEvaluator := &mockEvaluator{returns: evaluation}
		sut := newFeatureFlagEvaluatorMatcher(
			targetEvaluator,
			&mockValueOperatorMatcher{returns: true},
		)

		// when
		matches, err := sut.matches(request, context, condition)

		// then
		assert.Equal(t, true, matches)
		assert.Equal(t, nil, err)
		assert.Equal(t, 1, targetEvaluator.count)
	})
}

type mockEvaluator struct {
	returns interface{}
	count   int
}

func (m *mockEvaluator) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, error) {
	m.count++
	switch r := m.returns.(type) {
	case evaluator.Evaluation:
		return r, nil
	case error:
		return nil, r
	}
	return nil, nil
}

type mockValueOperatorMatcher struct {
	returns bool
}

func (m *mockValueOperatorMatcher) Matches(userValue interface{}, match model.TargetMatch) bool {
	return m.returns
}
