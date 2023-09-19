package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/bucketer"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/flow"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/target"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type EvaluationFlowFactory interface {
	Get(experimentType model.ExperimentType) (flow.EvaluationFlow, error)
}

func NewFlowFactory(targetMatcher target.Matcher, bucketer bucketer.Bucketer) EvaluationFlowFactory {

	actionResolver := &actionResolver{bucketer: bucketer}
	overrideResolver := &overrideResolver{targetMatcher, actionResolver}
	containerResolver := &containerResolver{bucketer: bucketer}
	targetDeterminer := &targetDeterminer{targetMatcher}
	targetRuleDeterminer := &targetRuleDeterminer{targetMatcher}

	abTestFlow := flow.NewEvaluationFlow(
		NewOverrideEvaluator(overrideResolver),
		NewIdentifierEvaluator(),
		NewContainerEvaluator(containerResolver),
		NewTargetEvaluator(targetDeterminer),
		NewDraftEvaluator(),
		NewPausedEvaluator(),
		NewCompletedEvaluator(),
		NewTrafficAllocatedEvaluator(actionResolver),
	)

	featureFlagFlow := flow.NewEvaluationFlow(
		NewDraftEvaluator(),
		NewPausedEvaluator(),
		NewCompletedEvaluator(),
		NewOverrideEvaluator(overrideResolver),
		NewIdentifierEvaluator(),
		NewTargetRuleEvaluator(targetRuleDeterminer, actionResolver),
		NewDefaultRuleEvaluator(actionResolver),
	)

	return &evaluationFlowFactory{
		abTestFlow:      abTestFlow,
		featureFlagFlow: featureFlagFlow,
	}
}

type evaluationFlowFactory struct {
	abTestFlow      flow.EvaluationFlow
	featureFlagFlow flow.EvaluationFlow
}

func (f *evaluationFlowFactory) Get(experimentType model.ExperimentType) (flow.EvaluationFlow, error) {
	switch experimentType {
	case model.ExperimentTypeAbTest:
		return f.abTestFlow, nil
	case model.ExperimentTypeFeatureFlag:
		return f.featureFlagFlow, nil
	}
	return nil, fmt.Errorf("unsupported experiment type [%s]", experimentType)
}
