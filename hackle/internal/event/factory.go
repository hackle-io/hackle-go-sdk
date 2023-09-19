package event

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/remoteconfig"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"
)

type Factory interface {
	Create(request evaluator.Request, evaluation evaluator.Evaluation) ([]UserEvent, error)
}

func NewFactory(clock clock.Clock) Factory {
	return &factory{
		clock: clock,
	}
}

type factory struct {
	clock clock.Clock
}

func (f *factory) Create(request evaluator.Request, evaluation evaluator.Evaluation) ([]UserEvent, error) {
	timestamp := f.clock.CurrentMillis()

	events := make([]UserEvent, 0)

	rootEvent, err := f.create(request, evaluation, timestamp, properties.NewBuilder())
	if err != nil {
		return nil, err
	}
	events = append(events, rootEvent)

	for _, targetEvaluation := range evaluation.TargetEvaluations() {
		builder := properties.NewBuilder()
		builder.Add(rootType, request.Key().Type.String())
		builder.Add(rootID, request.Key().ID)
		targetEvent, err := f.create(request, targetEvaluation, timestamp, builder)
		if err != nil {
			return nil, err
		}
		events = append(events, targetEvent)
	}
	return events, nil
}

func (f *factory) create(
	request evaluator.Request,
	evaluation evaluator.Evaluation,
	timestamp int64,
	properties *properties.Builder,
) (UserEvent, error) {
	switch e := evaluation.(type) {
	case experiment.Evaluation:
		parameterConfigID := e.ParameterConfigID()
		if parameterConfigID != nil {
			properties.Add(configIdPropertyKey, *parameterConfigID)
		}
		properties.Add(experimentVersionKey, e.Experiment.Version)
		properties.Add(executionVersionKey, e.Experiment.ExecutionVersion)
		return NewExposureEvent(e, properties.Build(), request.User(), timestamp), nil
	case remoteconfig.Evaluation:
		properties.AddAll(e.Properties)
		return NewRemoteConfigEvent(e, properties.Build(), request.User(), timestamp), nil
	}

	return nil, fmt.Errorf("unsupported evaluator.Evaluation [%T]", evaluation)
}

const (
	rootType = "$targetingRootType"
	rootID   = "$targetingRootId"

	configIdPropertyKey = "$parameterConfigurationId"

	experimentVersionKey = "$experiment_version"
	executionVersionKey  = "$execution_version"
)
