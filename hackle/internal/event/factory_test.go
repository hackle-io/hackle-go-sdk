package event

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/remoteconfig"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFactory_Create(t *testing.T) {

	t.Run("create", func(t *testing.T) {
		sut := NewFactory(clock.Fixed(42))

		context := evaluator.NewContext()

		evaluation1 := experiment.NewEvaluationOf(
			decision.ReasonTrafficAllocated,
			make([]evaluator.Evaluation, 0),
			model.Experiment{ID: 1, Type: model.ExperimentTypeAbTest, Version: 1, ExecutionVersion: 1},
			ref.Int64(42),
			"B",
			&model.ParameterConfiguration{ID: 42},
		)
		evaluation2 := experiment.NewEvaluationOf(
			decision.ReasonDefaultRule,
			make([]evaluator.Evaluation, 0),
			model.Experiment{ID: 2, Type: model.ExperimentTypeFeatureFlag, Version: 2, ExecutionVersion: 3},
			ref.Int64(320),
			"A",
			nil,
		)

		context.AddEvaluation(evaluation1)
		context.AddEvaluation(evaluation2)

		hackleUser := user.NewHackleUserBuilder().Identifier("a", "b").Build()
		request := remoteconfig.NewRequest(
			mocks.CreateWorkspace(),
			hackleUser,
			model.RemoteConfigParameter{ID: 2000},
			types.String,
			"default",
		)
		evaluation := remoteconfig.NewEvaluation(
			request,
			context,
			ref.Int64(999),
			"RC",
			decision.ReasonTargetRuleMatch,
			properties.NewBuilder(),
		)

		events, err := sut.Create(request, evaluation)
		assert.Equal(t, nil, err)
		assert.Equal(t, 3, len(events))

		event0 := events[0].(RemoteConfigEvent)
		assert.Equal(t, int64(42), event0.Timestamp())
		assert.Equal(t, hackleUser, event0.User())
		assert.Equal(t, model.RemoteConfigParameter{ID: 2000}, event0.Parameter)
		assert.Equal(t, ref.Int64(999), event0.ValueID)
		assert.Equal(t, "TARGET_RULE_MATCH", event0.DecisionReason)
		assert.Equal(t, map[string]interface{}{
			"returnValue": "RC",
		}, event0.Properties)

		event1 := events[1].(ExposureEvent)
		assert.Equal(t, int64(42), event1.Timestamp())
		assert.Equal(t, hackleUser, event1.User())
		assert.Equal(t, evaluation1.Experiment, event1.Experiment)
		assert.Equal(t, ref.Int64(42), event1.VariationID)
		assert.Equal(t, "B", event1.VariationKey)
		assert.Equal(t, "TRAFFIC_ALLOCATED", event1.DecisionReason)
		assert.Equal(t, map[string]interface{}{
			"$targetingRootType":        "REMOTE_CONFIG",
			"$targetingRootId":          int64(2000),
			"$parameterConfigurationId": int64(42),
			"$experiment_version":       1,
			"$execution_version":        1,
		}, event1.Properties)

		event2 := events[2].(ExposureEvent)
		assert.Equal(t, int64(42), event1.Timestamp())
		assert.Equal(t, hackleUser, event2.User())
		assert.Equal(t, evaluation2.Experiment, event2.Experiment)
		assert.Equal(t, ref.Int64(320), event2.VariationID)
		assert.Equal(t, "A", event2.VariationKey)
		assert.Equal(t, "DEFAULT_RULE", event2.DecisionReason)
		assert.Equal(t, map[string]interface{}{
			"$targetingRootType":  "REMOTE_CONFIG",
			"$targetingRootId":    int64(2000),
			"$experiment_version": 2,
			"$execution_version":  3,
		}, event2.Properties)
	})

	t.Run("unsupported", func(t *testing.T) {
		sut := NewFactory(clock.Fixed(42))

		_, err := sut.Create(evaluator.SimpleRequest{}, evaluator.SimpleEvaluation{})
		assert.Equal(t, errors.New("unsupported evaluator.Evaluation [evaluator.SimpleEvaluation]"), err)

		context := evaluator.NewContext()
		context.AddEvaluation(evaluator.SimpleEvaluation{})
		request := remoteconfig.NewRequest(
			mocks.CreateWorkspace(),
			user.HackleUser{},
			model.RemoteConfigParameter{ID: 2000},
			types.String,
			"default",
		)
		evaluation := remoteconfig.NewEvaluation(
			request,
			context,
			ref.Int64(999),
			"RC",
			decision.ReasonTargetRuleMatch,
			properties.NewBuilder(),
		)

		_, err = sut.Create(request, evaluation)
		assert.Equal(t, errors.New("unsupported evaluator.Evaluation [evaluator.SimpleEvaluation]"), err)
	})
}
