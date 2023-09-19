package experiment

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/flow"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEvaluator(t *testing.T) {
	e := NewEvaluator(&mockEvaluationFlowFactory{})
	ee := e.(*experimentEvaluator)
	assert.NotNil(t, ee.Evaluator)
	assert.Equal(t, ee, ee.Evaluator)
	assert.Equal(t, ee, ee.BaseEvaluator.Evaluator)
	assert.NotNil(t, ee.flowFactory)
}

func TestExperimentEvaluator_Supports(t *testing.T) {
	sut := experimentEvaluator{}

	assert.True(t, sut.Supports(Request{}))
	assert.False(t, sut.Supports(evaluator.SimpleRequest{}))
}

func TestExperimentEvaluator_Evaluate(t *testing.T) {

	t.Run("when request is not experiment request type then return error", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{}

		sut := experimentEvaluator{}

		// when
		_, err := sut.Evaluate(request, evaluator.NewContext())

		// then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unsupported evaluator request")
	})

	t.Run("EvaluateExperiment", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{Type: model.ExperimentTypeAbTest}, DefaultVariationKey: "A"}
		evaluation := Evaluation{reason: "42"}
		factory := &mockEvaluationFlowFactory{flow: &mockFlow{returns: evaluation}}
		sut := experimentEvaluator{flowFactory: factory}

		// when
		actual, err := sut.Evaluate(request, evaluator.NewContext())

		// then
		assert.Nil(t, err)
		assert.Equal(t, evaluation, actual)
	})
}

func TestExperimentEvaluator_EvaluateExperiment(t *testing.T) {
	t.Run("when cannot get flow then return error", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{Type: model.ExperimentTypeAbTest}}

		sut := experimentEvaluator{flowFactory: &mockEvaluationFlowFactory{}}

		// when
		_, err := sut.EvaluateExperiment(request, evaluator.NewContext())

		// then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "flow not found")
	})

	t.Run("when error on flow evaluate then return error", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{Type: model.ExperimentTypeAbTest}}
		factory := &mockEvaluationFlowFactory{flow: &mockFlow{returns: errors.New("flow error")}}
		sut := experimentEvaluator{flowFactory: factory}

		// when
		_, err := sut.EvaluateExperiment(request, evaluator.NewContext())

		// then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "flow error")
	})

	t.Run("when flow evaluated as nil then return default evaluation", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{Type: model.ExperimentTypeAbTest}, DefaultVariationKey: "A"}
		factory := &mockEvaluationFlowFactory{flow: &mockFlow{returns: nil}}
		sut := experimentEvaluator{flowFactory: factory}

		// when
		actual, err := sut.EvaluateExperiment(request, evaluator.NewContext())

		// then
		assert.Nil(t, err)
		assert.Equal(t, "A", actual.VariationKey)
		assert.Equal(t, "TRAFFIC_NOT_ALLOCATED", actual.Reason())
	})

	t.Run("when flow evaluated as not experiment evaluation then return error", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{Type: model.ExperimentTypeAbTest}, DefaultVariationKey: "A"}
		factory := &mockEvaluationFlowFactory{flow: &mockFlow{returns: evaluator.SimpleEvaluation{}}}
		sut := experimentEvaluator{flowFactory: factory}

		// when
		_, err := sut.EvaluateExperiment(request, evaluator.NewContext())

		// then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unexpected evaluation")
	})

	t.Run("when flow evaluated as experiment evaluation then return evaluated evaluation", func(t *testing.T) {
		// given
		request := Request{Experiment: model.Experiment{Type: model.ExperimentTypeAbTest}, DefaultVariationKey: "A"}
		evaluation := Evaluation{reason: "42"}
		factory := &mockEvaluationFlowFactory{flow: &mockFlow{returns: evaluation}}
		sut := experimentEvaluator{flowFactory: factory}

		// when
		actual, err := sut.EvaluateExperiment(request, evaluator.NewContext())

		// then
		assert.Nil(t, err)
		assert.Equal(t, evaluation, actual)
	})
}

type mockEvaluationFlowFactory struct {
	flow flow.EvaluationFlow
}

func (m *mockEvaluationFlowFactory) Get(experimentType model.ExperimentType) (flow.EvaluationFlow, error) {
	if m.flow == nil {
		return nil, errors.New("flow not found")
	}
	return m.flow, nil
}
