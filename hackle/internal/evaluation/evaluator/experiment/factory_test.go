package experiment

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/flow"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvaluationFlowFactory_Get(t *testing.T) {

	decisionWith := func(f flow.EvaluationFlow, e interface{}) flow.EvaluationFlow {
		decision, ok := f.(*flow.Decision)
		assert.True(t, ok)
		assert.IsType(t, e, decision.Evaluator)
		return decision.NextFlow
	}

	end := func(f flow.EvaluationFlow) {
		assert.IsType(t, f, &flow.End{})
	}

	t.Run("AB_TEST", func(t *testing.T) {
		factory := NewFlowFactory(nil, nil)
		f, err := factory.Get(model.ExperimentTypeAbTest)
		assert.Nil(t, err)
		f = decisionWith(f, &OverrideEvaluator{})
		f = decisionWith(f, &IdentifierEvaluator{})
		f = decisionWith(f, &ContainerEvaluator{})
		f = decisionWith(f, &TargetEvaluator{})
		f = decisionWith(f, &DraftEvaluator{})
		f = decisionWith(f, &PausedEvaluator{})
		f = decisionWith(f, &CompletedEvaluator{})
		f = decisionWith(f, &TrafficAllocateEvaluator{})
		end(f)
	})

	t.Run("FEATURE_FLAG", func(t *testing.T) {
		factory := NewFlowFactory(nil, nil)
		f, err := factory.Get(model.ExperimentTypeFeatureFlag)
		assert.Nil(t, err)
		f = decisionWith(f, &DraftEvaluator{})
		f = decisionWith(f, &PausedEvaluator{})
		f = decisionWith(f, &CompletedEvaluator{})
		f = decisionWith(f, &OverrideEvaluator{})
		f = decisionWith(f, &IdentifierEvaluator{})
		f = decisionWith(f, &TargetRuleEvaluator{})
		f = decisionWith(f, &DefaultRuleEvaluator{})
		end(f)
	})

	t.Run("unsupported type", func(t *testing.T) {
		factory := NewFlowFactory(nil, nil)
		_, err := factory.Get("INVALID")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unsupported experiment type")
	})
}
