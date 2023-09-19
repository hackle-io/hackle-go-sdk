package experiment

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/flow"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseFlowEvaluator_Evaluate(t *testing.T) {

	t.Run("when request is not experiment type then return error", func(t *testing.T) {
		// given
		sut := &baseFlowEvaluator{}

		// when
		_, _, err := sut.Evaluate(evaluator.SimpleRequest{}, evaluator.NewContext(), flow.NewEvaluationFlow())

		// then
		assert.NotNil(t, err)
	})

	t.Run("when error on flow evaluate then return error", func(t *testing.T) {
		// given
		sut := &baseFlowEvaluator{&mockFlowEvaluator{errors.New("failed to evaluate")}}

		// when
		_, _, err := sut.Evaluate(Request{}, evaluator.NewContext(), flow.NewEvaluationFlow())

		// then
		assert.Equal(t, errors.New("failed to evaluate"), err)
	})

	t.Run("when flow evaluated then return evaluated evaluation", func(t *testing.T) {
		// given
		evaluation := Evaluation{reason: "42"}
		sut := &baseFlowEvaluator{&mockFlowEvaluator{evaluation}}

		// when
		actual, ok, err := sut.Evaluate(Request{}, evaluator.NewContext(), flow.NewEvaluationFlow())

		// then
		assert.Equal(t, evaluation, actual)
		assert.True(t, ok)
		assert.Nil(t, err)
	})
}

func TestBaseFlowEvaluator_evaluation(t *testing.T) {
	t.Run("when error on new evaluation then return error", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:         42,
			Type:       model.ExperimentTypeAbTest,
			Status:     model.ExperimentStatusRunning,
			Variations: []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			workspace:           mocks.CreateWorkspace(),
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		variation := model.Variation{ID: 1002, Key: "B", ParameterConfigurationID: ref.Int64(320)}
		reason := decision.ReasonTrafficAllocated

		// when
		sut := &baseFlowEvaluator{}
		_, ok, err := sut.evaluation(request, context, variation, reason)

		// then
		assert.Equal(t, false, ok)
		assert.NotNil(t, err)
	})
}

func TestBaseEvaluator_evaluateDefault(t *testing.T) {
	t.Run("when error on new evaluation default then return error", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:         42,
			Type:       model.ExperimentTypeAbTest,
			Status:     model.ExperimentStatusRunning,
			Variations: []model.Variation{{ID: 1001, Key: "A", ParameterConfigurationID: ref.Int64(320)}, {ID: 1002, Key: "B"}},
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			workspace:           mocks.CreateWorkspace(),
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		reason := decision.ReasonTrafficAllocated

		// when
		sut := &baseFlowEvaluator{}
		_, ok, err := sut.defaultEvaluation(request, context, reason)

		// then
		assert.Equal(t, false, ok)
		assert.NotNil(t, err)
	})
}

func TestOverrideEvaluator_evaluate(t *testing.T) {

	t.Run("when error on override resolve then return error", func(t *testing.T) {
		// given
		resolver := &mockOverrideResolver{returns: errors.New("failed to resolve")}

		request := Request{Experiment: model.Experiment{ID: 42, Type: model.ExperimentTypeAbTest}}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewOverrideEvaluator(resolver)
		evaluation, ok, err := sut.Evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to resolve")
	})

	t.Run("when resolve as nil then evaluate next flow", func(t *testing.T) {
		// given
		resolver := &mockOverrideResolver{returns: nil}

		request := Request{Experiment: model.Experiment{ID: 42, Type: model.ExperimentTypeAbTest}}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewOverrideEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Nil(t, err)
	})

	t.Run("when ab test overridden then return overridden variation with overridden reason", func(t *testing.T) {
		// given
		resolver := &mockOverrideResolver{returns: model.Variation{ID: 320, Key: "B"}}

		request := Request{Experiment: model.Experiment{ID: 42, Type: model.ExperimentTypeAbTest}}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewOverrideEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "OVERRIDDEN",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        model.Experiment{ID: 42, Type: model.ExperimentTypeAbTest},
			VariationID:       ref.Int64(320),
			VariationKey:      "B",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Nil(t, err)
	})

	t.Run("when feature flag overridden then return overridden variation with individual match reason", func(t *testing.T) {
		// given
		resolver := &mockOverrideResolver{returns: model.Variation{ID: 320, Key: "B"}}

		request := Request{Experiment: model.Experiment{ID: 42, Type: model.ExperimentTypeFeatureFlag}}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewOverrideEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "INDIVIDUAL_TARGET_MATCH",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        model.Experiment{ID: 42, Type: model.ExperimentTypeFeatureFlag},
			VariationID:       ref.Int64(320),
			VariationKey:      "B",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Nil(t, err)
	})

	t.Run("when unsupported experiment type then return error", func(t *testing.T) {
		// given
		resolver := &mockOverrideResolver{returns: model.Variation{ID: 320, Key: "B"}}

		request := Request{Experiment: model.Experiment{ID: 42, Type: model.ExperimentType("unsupported")}}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		sut := NewOverrideEvaluator(resolver)

		// when
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unsupported experiment type")
	})
}
func TestDraftEvaluator_evaluate(t *testing.T) {

	t.Run("when draft experiment then return default variation", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusDraft,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDraftEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "EXPERIMENT_DRAFT",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       nil,
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when not draft experiment then evaluate next flow", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusRunning,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDraftEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestPausedEvaluator_evaluate(t *testing.T) {

	t.Run("when paused ab test then return default variation", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusPaused,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewPausedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "EXPERIMENT_PAUSED",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       nil,
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when paused feature flag then return default variation", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeFeatureFlag,
			Status: model.ExperimentStatusPaused,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewPausedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "FEATURE_FLAG_INACTIVE",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       nil,
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when not paused experiment then evaluate next flow", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusRunning,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewPausedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when unsupported experiment type then return error", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentType("unsupported"),
			Status: model.ExperimentStatusPaused,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewPausedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("unsupported experiment type [unsupported]"), err)
	})
}

func TestCompletedEvaluator_evaluate(t *testing.T) {
	t.Run("when completed experiment then return winner variation", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusCompleted,
			Variations: []model.Variation{
				{ID: 1001, Key: "A"},
				{ID: 1002, Key: "B"},
			},
			WinnerVariationID: ref.Int64(1002),
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewCompletedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "EXPERIMENT_COMPLETED",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1002),
			VariationKey:      "B",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when completed experiment without winner then return error", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusCompleted,
			Variations: []model.Variation{
				{ID: 1001, Key: "A"},
				{ID: 1002, Key: "B"},
			},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewCompletedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("winner variation [42]"), err)
	})

	t.Run("when not completed experiment then evaluate next flow", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:     42,
			Type:   model.ExperimentTypeAbTest,
			Status: model.ExperimentStatusRunning,
			Variations: []model.Variation{
				{ID: 1001, Key: "A"},
				{ID: 1002, Key: "B"},
			},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewCompletedEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestTargetEvaluator_evaluate(t *testing.T) {
	t.Run("when experiment is not ab test type then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetDeterminer{returns: false}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeFeatureFlag,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetEvaluator(determiner)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment type must be AB_TEST [42]"), err)
	})

	t.Run("when error on determine then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetDeterminer{returns: errors.New("determine error")}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetEvaluator(determiner)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("determine error"), err)
	})

	t.Run("when user in experiment target then evaluate next flow", func(t *testing.T) {
		// given
		determiner := &mockTargetDeterminer{returns: true}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetEvaluator(determiner)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when user not in experiment target then return default variation", func(t *testing.T) {
		// given
		determiner := &mockTargetDeterminer{returns: false}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetEvaluator(determiner)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "NOT_IN_EXPERIMENT_TARGET",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1001),
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestTrafficAllocateEvaluator_evaluate(t *testing.T) {
	t.Run("when experiment is not running then return error", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusDraft,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTrafficAllocatedEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment status must be RUNNING [42]"), err)
	})

	t.Run("when experiment is not ab test type then return error", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeFeatureFlag,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTrafficAllocatedEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment type must be AB_TEST [42]"), err)
	})

	t.Run("when error on resolve variation then return error", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: errors.New("resolve error")}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTrafficAllocatedEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("resolve error"), err)
	})

	t.Run("when cannot resolve variation then return default variation with traffic not allocated", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTrafficAllocatedEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "TRAFFIC_NOT_ALLOCATED",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1001),
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when resolved variation is dropped then return default variation", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: model.Variation{ID: 1002, Key: "B", IsDropped: true}}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTrafficAllocatedEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "VARIATION_DROPPED",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1001),
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when variation decided then return that variation", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: model.Variation{ID: 1002, Key: "B", IsDropped: false}}

		experiment := model.Experiment{
			ID:                42,
			Type:              model.ExperimentTypeAbTest,
			Status:            model.ExperimentStatusRunning,
			Variations:        []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			WinnerVariationID: nil,
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTrafficAllocatedEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "TRAFFIC_ALLOCATED",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1002),
			VariationKey:      "B",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestNewTargetRuleEvaluator_evaluate(t *testing.T) {
	t.Run("when experiment is not running then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: nil}
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:         42,
			Type:       model.ExperimentTypeFeatureFlag,
			Status:     model.ExperimentStatusDraft,
			Variations: []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment status must be RUNNING [42]"), err)
	})

	t.Run("when experiment is not feature flag type then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: nil}
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:         42,
			Type:       model.ExperimentTypeAbTest,
			Status:     model.ExperimentStatusRunning,
			Variations: []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment type must be FEATURE_FLAG [42]"), err)
	})

	t.Run("when identifier not exist then evaluate next flow", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: nil}
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "custom_id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when error on determine then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: errors.New("determine error")}
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("determine error"), err)
	})

	t.Run("when cannot determine then evaluate next flow", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: nil}
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when error on resolve variation then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: model.TargetRule{}}
		resolver := &mockActionResolver{returns: errors.New("resolve error")}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("resolve error"), err)
	})

	t.Run("when cannot resolve then return error", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: model.TargetRule{}}
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("feature flag must decide variation [42]"), err)
	})

	t.Run("when variation resolved then return resolved variation", func(t *testing.T) {
		// given
		determiner := &mockTargetRuleDeterminer{returns: model.TargetRule{}}
		resolver := &mockActionResolver{returns: model.Variation{ID: 1002, Key: "B"}}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewTargetRuleEvaluator(determiner, resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "TARGET_RULE_MATCH",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1002),
			VariationKey:      "B",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestDefaultRuleEvaluator_evaluate(t *testing.T) {
	t.Run("when experiment is not running then return error", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:         42,
			Type:       model.ExperimentTypeFeatureFlag,
			Status:     model.ExperimentStatusDraft,
			Variations: []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDefaultRuleEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment status must be RUNNING [42]"), err)
	})

	t.Run("when experiment is not feature flag type then return error", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:         42,
			Type:       model.ExperimentTypeAbTest,
			Status:     model.ExperimentStatusRunning,
			Variations: []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
		}
		request := Request{
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDefaultRuleEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("experiment type must be FEATURE_FLAG [42]"), err)
	})

	t.Run("when identifier not found then return default variation", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "custom_id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDefaultRuleEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "DEFAULT_RULE",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1001),
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when error on resolve variation then return default variation", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: errors.New("resolve error")}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDefaultRuleEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("resolve error"), err)
	})

	t.Run("when cannot resolve variation then return error", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: nil}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDefaultRuleEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("feature flag must decide variation [42]"), err)
	})
	t.Run("when variation decided then return that variation", func(t *testing.T) {
		// given
		resolver := &mockActionResolver{returns: model.Variation{ID: 1002, Key: "B"}}

		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeFeatureFlag,
			Status:         model.ExperimentStatusRunning,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "$id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewDefaultRuleEvaluator(resolver)
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "DEFAULT_RULE",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1002),
			VariationKey:      "B",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestContainerEvaluator_evaluate(t *testing.T) {

	t.Run("when not mutually exclusion experiment then evaluate next flow", func(t *testing.T) {
		// given
		resolver := &mockContainerResolver{returns: false}

		experiment := model.Experiment{
			ID:          42,
			Type:        model.ExperimentTypeAbTest,
			Status:      model.ExperimentStatusRunning,
			Variations:  []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			ContainerID: nil,
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewContainerEvaluator(resolver)

		// then
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when container not found then return error", func(t *testing.T) {
		// given
		resolver := &mockContainerResolver{returns: false}

		experiment := model.Experiment{
			ID:          42,
			Type:        model.ExperimentTypeAbTest,
			Status:      model.ExperimentStatusRunning,
			Variations:  []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			ContainerID: ref.Int64(320),
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			workspace:           mocks.CreateWorkspace(),
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewContainerEvaluator(resolver)

		// then
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("container [320]"), err)
	})

	t.Run("when error on resolve container group then return error", func(t *testing.T) {
		// given
		resolver := &mockContainerResolver{returns: errors.New("container error")}

		container := model.Container{ID: 320}
		experiment := model.Experiment{
			ID:          42,
			Type:        model.ExperimentTypeAbTest,
			Status:      model.ExperimentStatusRunning,
			Variations:  []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			ContainerID: ref.Int64(320),
		}
		workspace := mocks.CreateWorkspace()
		workspace.Container(container)
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			workspace:           workspace,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewContainerEvaluator(resolver)

		// then
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, nil, evaluation)
		assert.Equal(t, false, ok)
		assert.Equal(t, errors.New("container error"), err)
	})

	t.Run("when user is in container group then evaluate next flow", func(t *testing.T) {
		// given
		resolver := &mockContainerResolver{returns: true}

		container := model.Container{ID: 320}
		experiment := model.Experiment{
			ID:          42,
			Type:        model.ExperimentTypeAbTest,
			Status:      model.ExperimentStatusRunning,
			Variations:  []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			ContainerID: ref.Int64(320),
		}
		workspace := mocks.CreateWorkspace()
		workspace.Container(container)
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			workspace:           workspace,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewContainerEvaluator(resolver)

		// then
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})

	t.Run("when user is not in container group then evaluate next flow", func(t *testing.T) {
		// given
		resolver := &mockContainerResolver{returns: false}

		container := model.Container{ID: 320}
		experiment := model.Experiment{
			ID:          42,
			Type:        model.ExperimentTypeAbTest,
			Status:      model.ExperimentStatusRunning,
			Variations:  []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			ContainerID: ref.Int64(320),
		}
		workspace := mocks.CreateWorkspace()
		workspace.Container(container)
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			workspace:           workspace,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewContainerEvaluator(resolver)

		// then
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "NOT_IN_MUTUAL_EXCLUSION_EXPERIMENT",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1001),
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

func TestIdentifierEvaluator_evaluate(t *testing.T) {
	t.Run("when identifier exist then evaluate next flow", func(t *testing.T) {
		// given
		request := Request{
			user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment: model.Experiment{
				ID:             42,
				Type:           model.ExperimentTypeAbTest,
				Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
				IdentifierType: "$id",
			},
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewIdentifierEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{reason: "next_flow"}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
	t.Run("when identifier not found then return default variation", func(t *testing.T) {
		// given
		experiment := model.Experiment{
			ID:             42,
			Type:           model.ExperimentTypeAbTest,
			Variations:     []model.Variation{{ID: 1001, Key: "A"}, {ID: 1002, Key: "B"}},
			IdentifierType: "custom_id",
		}
		request := Request{
			user:                user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
			Experiment:          experiment,
			DefaultVariationKey: "A",
		}
		context := evaluator.NewContext()
		nextFlow := &mockFlow{returns: Evaluation{reason: "next_flow"}}

		// when
		sut := NewIdentifierEvaluator()
		evaluation, ok, err := sut.evaluate(request, context, nextFlow)

		// then
		assert.Equal(t, Evaluation{
			reason:            "IDENTIFIER_NOT_FOUND",
			targetEvaluations: make([]evaluator.Evaluation, 0),
			Experiment:        experiment,
			VariationID:       ref.Int64(1001),
			VariationKey:      "A",
			config:            nil,
		}, evaluation)
		assert.Equal(t, true, ok)
		assert.Equal(t, nil, err)
	})
}

type mockFlow struct {
	returns interface{}
	count   int
}

func (m *mockFlow) Evaluate(request evaluator.Request, context evaluator.Context) (evaluator.Evaluation, bool, error) {
	m.count++
	switch r := m.returns.(type) {
	case evaluator.Evaluation:
		return r, true, nil
	case error:
		return nil, false, r
	default:
		return nil, false, nil
	}
}

type mockFlowEvaluator struct {
	returns interface{}
}

func (m *mockFlowEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	switch r := m.returns.(type) {
	case evaluator.Evaluation:
		return r, true, nil
	case error:
		return nil, false, r
	default:
		return nil, false, nil
	}
}

type mockOverrideResolver struct {
	returns interface{}
}

func (m *mockOverrideResolver) Resolve(request Request, context evaluator.Context) (model.Variation, bool, error) {
	switch r := m.returns.(type) {
	case model.Variation:
		return r, true, nil
	case error:
		return model.Variation{}, false, r
	default:
		return model.Variation{}, false, nil
	}
}

type mockTargetDeterminer struct {
	returns interface{}
}

func (m *mockTargetDeterminer) IsUserInExperimentTarget(request Request, context evaluator.Context) (bool, error) {
	switch r := m.returns.(type) {
	case bool:
		return r, nil
	case error:
		return false, r
	default:
		return false, nil
	}
}

type mockActionResolver struct {
	returns interface{}
	count   int
}

func (m *mockActionResolver) Resolve(request Request, action model.Action) (model.Variation, bool, error) {
	m.count++
	switch r := m.returns.(type) {
	case model.Variation:
		return r, true, nil
	case error:
		return model.Variation{}, false, r
	default:
		return model.Variation{}, false, nil
	}
}

type mockTargetRuleDeterminer struct {
	returns interface{}
}

func (m *mockTargetRuleDeterminer) Determine(request Request, context evaluator.Context) (model.TargetRule, bool, error) {
	switch r := m.returns.(type) {
	case model.TargetRule:
		return r, true, nil
	case error:
		return model.TargetRule{}, false, r
	default:
		return model.TargetRule{}, false, nil
	}
}

type mockContainerResolver struct {
	returns interface{}
}

func (m *mockContainerResolver) IsUserInContainerGroup(request Request, container model.Container) (bool, error) {
	switch r := m.returns.(type) {
	case bool:
		return r, nil
	case error:
		return false, r
	default:
		return false, nil
	}
}
