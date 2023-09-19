package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
)

type Evaluation struct {
	reason            string
	targetEvaluations []evaluator.Evaluation
	Experiment        model.Experiment
	VariationID       *int64
	VariationKey      string
	config            *model.ParameterConfiguration
}

func NewEvaluation(
	request Request,
	context evaluator.Context,
	variation model.Variation,
	reason string,
) (Evaluation, error) {
	c, _, err := configuration(request.Workspace(), variation)
	if err != nil {
		return Evaluation{}, err
	}
	return Evaluation{
		reason:            reason,
		targetEvaluations: context.Evaluations(),
		Experiment:        request.Experiment,
		VariationID:       &variation.ID,
		VariationKey:      variation.Key,
		config:            c,
	}, nil
}

func NewEvaluationDefault(
	request Request,
	context evaluator.Context,
	reason string,
) (Evaluation, error) {
	variation, ok := request.Experiment.GetVariationByKey(request.DefaultVariationKey)
	if !ok {
		return Evaluation{
			reason:            reason,
			targetEvaluations: context.Evaluations(),
			Experiment:        request.Experiment,
			VariationID:       nil,
			VariationKey:      request.DefaultVariationKey,
			config:            nil,
		}, nil
	}

	return NewEvaluation(request, context, variation, reason)
}

func NewEvaluationOf(
	reason string,
	targetEvaluations []evaluator.Evaluation,
	experiment model.Experiment,
	variationID *int64,
	variationKey string,
	config *model.ParameterConfiguration,
) Evaluation {
	return Evaluation{
		reason:            reason,
		targetEvaluations: targetEvaluations,
		Experiment:        experiment,
		VariationID:       variationID,
		VariationKey:      variationKey,
		config:            config,
	}
}

func (e Evaluation) Reason() string {
	return e.reason
}

func (e Evaluation) TargetEvaluations() []evaluator.Evaluation {
	return e.targetEvaluations
}

func (e Evaluation) String() string {
	return fmt.Sprintf("ExperimentEvaluation(experimentID=%d, experimentKey=%d, variation=%s, reason=%s)", e.Experiment.ID, e.Experiment.Key, e.VariationKey, e.reason)
}

func (e Evaluation) With(reason string) Evaluation {
	e.reason = reason
	return e
}
func (e Evaluation) Config() config.Config {
	if e.config == nil {
		return config.Empty()
	}
	return config.New(e.config.Parameters)
}

func (e Evaluation) ParameterConfigID() *int64 {
	if e.config == nil {
		return nil
	}
	return &e.config.ID
}

func configuration(workspace workspace.Workspace, variation model.Variation) (*model.ParameterConfiguration, bool, error) {
	parameterConfigurationID := variation.ParameterConfigurationID
	if parameterConfigurationID == nil {
		return nil, false, nil
	}
	parameterConfiguration, ok := workspace.GetParameterConfiguration(*parameterConfigurationID)
	if !ok {
		return nil, false, fmt.Errorf("parameterConfiguration[%d]", parameterConfigurationID)
	}
	return &parameterConfiguration, true, nil
}
