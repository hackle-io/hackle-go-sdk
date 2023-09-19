package core

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/remoteconfig"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/event"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
)

type Core interface {
	Experiment(experimentKey int64, user user.HackleUser, defaultVariation string) (decision.ExperimentDecision, error)
	FeatureFlag(featureKey int64, user user.HackleUser) (decision.FeatureFlagDecision, error)
	RemoteConfig(parameterKey string, user user.HackleUser, requiredType types.ValueType, defaultValue interface{}) (decision.RemoteConfigDecision, error)
	Track(e event.HackleEvent, user user.HackleUser)
	Close()
}

func New(workspaceFetcher workspace.Fetcher, eventProcessor event.Processor) Core {
	experimentEvaluator, remoteConfigEvaluator := evaluation.NewEvaluators()
	return &core{
		experimentEvaluator:   experimentEvaluator,
		remoteConfigEvaluator: remoteConfigEvaluator,
		workspaceFetcher:      workspaceFetcher,
		eventFactory:          event.NewFactory(clock.System),
		eventProcessor:        eventProcessor,
		clock:                 clock.System,
	}
}

type core struct {
	experimentEvaluator   experiment.Evaluator
	remoteConfigEvaluator remoteconfig.Evaluator
	workspaceFetcher      workspace.Fetcher
	eventFactory          event.Factory
	eventProcessor        event.Processor
	clock                 clock.Clock
}

func (c *core) Experiment(experimentKey int64, user user.HackleUser, defaultVariation string) (decision.ExperimentDecision, error) {
	ws, ok := c.workspaceFetcher.Fetch()
	if !ok {
		return decision.NewExperimentDecision(defaultVariation, decision.ReasonSdkNotReady, config.Empty()), nil
	}

	exp, ok := ws.GetExperiment(experimentKey)
	if !ok {
		return decision.NewExperimentDecision(defaultVariation, decision.ReasonExperimentNotFound, config.Empty()), nil
	}

	req := experiment.NewRequest(ws, user, exp, defaultVariation)
	eval, err := c.experimentEvaluator.EvaluateExperiment(req, evaluator.NewContext())
	if err != nil {
		return decision.ExperimentDecision{}, err
	}

	events, err := c.eventFactory.Create(req, eval)
	if err != nil {
		return decision.ExperimentDecision{}, err
	}

	for _, it := range events {
		c.eventProcessor.Process(it)
	}

	return decision.NewExperimentDecision(eval.VariationKey, eval.Reason(), eval.Config()), nil
}

func (c *core) FeatureFlag(featureKey int64, user user.HackleUser) (decision.FeatureFlagDecision, error) {
	ws, ok := c.workspaceFetcher.Fetch()
	if !ok {
		return decision.NewFeatureFlagDecision(false, decision.ReasonSdkNotReady, config.Empty()), nil
	}

	flag, ok := ws.GetFeatureFlag(featureKey)
	if !ok {
		return decision.NewFeatureFlagDecision(false, decision.ReasonFeatureFlagNotFound, config.Empty()), nil
	}

	req := experiment.NewRequest(ws, user, flag, "A")
	eval, err := c.experimentEvaluator.EvaluateExperiment(req, evaluator.NewContext())
	if err != nil {
		return decision.FeatureFlagDecision{}, err
	}

	events, err := c.eventFactory.Create(req, eval)
	if err != nil {
		return decision.FeatureFlagDecision{}, err
	}

	for _, it := range events {
		c.eventProcessor.Process(it)
	}

	isOn := eval.VariationKey != "A"
	return decision.NewFeatureFlagDecision(isOn, eval.Reason(), eval.Config()), nil
}

func (c *core) RemoteConfig(parameterKey string, user user.HackleUser, requiredType types.ValueType, defaultValue interface{}) (decision.RemoteConfigDecision, error) {
	ws, ok := c.workspaceFetcher.Fetch()
	if !ok {
		return decision.NewRemoteConfigDecision(defaultValue, decision.ReasonSdkNotReady), nil
	}

	param, ok := ws.GetRemoteConfigParameter(parameterKey)
	if !ok {
		return decision.NewRemoteConfigDecision(defaultValue, decision.ReasonRemoteConfigParameterNotFound), nil
	}

	req := remoteconfig.NewRequest(ws, user, param, requiredType, defaultValue)
	eval, err := c.remoteConfigEvaluator.EvaluateRemoteConfig(req, evaluator.NewContext())
	if err != nil {
		return decision.RemoteConfigDecision{}, err
	}

	events, err := c.eventFactory.Create(req, eval)
	if err != nil {
		return decision.RemoteConfigDecision{}, err
	}
	for _, it := range events {
		c.eventProcessor.Process(it)
	}

	return decision.NewRemoteConfigDecision(eval.Value, eval.Reason()), nil
}

func (c *core) Track(e event.HackleEvent, user user.HackleUser) {
	eventType := c.eventType(e)
	trackEvent := event.NewTrackEvent(eventType, e, user, c.clock.CurrentMillis())
	c.eventProcessor.Process(trackEvent)
}

func (c *core) eventType(event event.HackleEvent) model.EventType {
	ws, ok := c.workspaceFetcher.Fetch()
	if !ok {
		return model.NewUndefinedEvent(event.Key())
	}
	if eventType, ok := ws.GetEventType(event.Key()); ok {
		return eventType
	} else {
		return model.NewUndefinedEvent(event.Key())
	}
}

func (c *core) Close() {
	c.eventProcessor.Close()
	c.workspaceFetcher.Close()
}
