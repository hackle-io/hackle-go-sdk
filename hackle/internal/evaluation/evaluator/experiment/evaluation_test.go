package experiment

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluation(t *testing.T) {

	evaluation1 := NewEvaluationOf(
		decision.ReasonTrafficAllocated,
		[]evaluator.Evaluation{Evaluation{}},
		model.Experiment{ID: 100},
		ref.Int64(42),
		"A",
		nil,
	)

	assert.Equal(t, "TRAFFIC_ALLOCATED", evaluation1.Reason())
	assert.Equal(t, []evaluator.Evaluation{Evaluation{}}, evaluation1.TargetEvaluations())
	assert.Equal(t, model.Experiment{ID: 100}, evaluation1.Experiment)
	assert.Equal(t, config.Empty(), evaluation1.Config())
	assert.Nil(t, evaluation1.ParameterConfigID())
	assert.Contains(t, evaluation1.String(), "ExperimentEvaluation")

	assert.Equal(t, "OVERRIDDEN", evaluation1.With(decision.ReasonOverridden).Reason())

	parameterConfiguration := model.ParameterConfiguration{
		ID:         1000,
		Parameters: map[string]interface{}{"a": "b"},
	}
	evaluation2, _ := NewEvaluation(
		Request{
			workspace: mocks.CreateWorkspace().ParameterConfiguration(parameterConfiguration),
			Experiment: model.Experiment{
				ID: 42,
			},
		},
		evaluator.NewContext(),
		model.Variation{ParameterConfigurationID: ref.Int64(1000)},
		decision.ReasonExperimentDraft,
	)

	assert.Equal(t, int64(1000), *evaluation2.ParameterConfigID())
	assert.Equal(t, config.New(map[string]interface{}{"a": "b"}), evaluation2.Config())
}
