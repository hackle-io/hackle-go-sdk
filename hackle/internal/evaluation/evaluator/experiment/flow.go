package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/flow"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type flowEvaluator interface {
	evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error)
}

type baseFlowEvaluator struct {
	flowEvaluator
}

func (e *baseFlowEvaluator) Evaluate(request evaluator.Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experimentRequest, ok := request.(Request)
	if !ok {
		return nil, false, fmt.Errorf("unsupported request: %T (expected: experiment.Request)", request)
	}
	experimentEvaluation, ok, err := e.evaluate(experimentRequest, context, nextFlow)
	if err != nil {
		return nil, false, err
	}
	return experimentEvaluation, ok, nil
}

func (e *baseFlowEvaluator) evaluation(request Request, context evaluator.Context, variation model.Variation, reason string) (Evaluation, bool, error) {
	evaluation, err := NewEvaluation(request, context, variation, reason)
	if err != nil {
		return Evaluation{}, false, err
	}
	return evaluation, true, nil
}

func (e *baseFlowEvaluator) defaultEvaluation(request Request, context evaluator.Context, reason string) (Evaluation, bool, error) {
	evaluation, err := NewEvaluationDefault(request, context, reason)
	if err != nil {
		return Evaluation{}, false, err
	}
	return evaluation, true, nil
}

type OverrideEvaluator struct {
	*baseFlowEvaluator
	overrideResolver OverrideResolver
}

func NewOverrideEvaluator(overrideResolver OverrideResolver) *OverrideEvaluator {
	e := &OverrideEvaluator{&baseFlowEvaluator{}, overrideResolver}
	e.flowEvaluator = e
	return e
}

func (e *OverrideEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	overriddenVariation, ok, err := e.overrideResolver.Resolve(request, context)
	if err != nil {
		return nil, false, err
	}
	if ok {
		switch request.Experiment.Type {
		case model.ExperimentTypeAbTest:
			return e.evaluation(request, context, overriddenVariation, decision.ReasonOverridden)
		case model.ExperimentTypeFeatureFlag:
			return e.evaluation(request, context, overriddenVariation, decision.ReasonIndividualTargetMatch)
		}
		return nil, false, fmt.Errorf("unsupported experiment type [%s]", request.Experiment.Type)
	}
	return nextFlow.Evaluate(request, context)
}

type DraftEvaluator struct {
	*baseFlowEvaluator
}

func NewDraftEvaluator() *DraftEvaluator {
	e := &DraftEvaluator{&baseFlowEvaluator{}}
	e.flowEvaluator = e
	return e
}

func (e *DraftEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	if request.Experiment.Status == model.ExperimentStatusDraft {
		return e.defaultEvaluation(request, context, decision.ReasonExperimentDraft)
	} else {
		return nextFlow.Evaluate(request, context)
	}
}

type PausedEvaluator struct {
	*baseFlowEvaluator
}

func NewPausedEvaluator() *PausedEvaluator {
	e := &PausedEvaluator{&baseFlowEvaluator{}}
	e.flowEvaluator = e
	return e
}

func (e *PausedEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experiment := request.Experiment
	if experiment.Status == model.ExperimentStatusPaused {
		switch experiment.Type {
		case model.ExperimentTypeAbTest:
			return e.defaultEvaluation(request, context, decision.ReasonExperimentPaused)
		case model.ExperimentTypeFeatureFlag:
			return e.defaultEvaluation(request, context, decision.ReasonFeatureFlagInactive)
		}
		return nil, false, fmt.Errorf("unsupported experiment type [%s]", experiment.Type)
	} else {
		return nextFlow.Evaluate(request, context)
	}
}

type CompletedEvaluator struct {
	*baseFlowEvaluator
}

func NewCompletedEvaluator() *CompletedEvaluator {
	e := &CompletedEvaluator{&baseFlowEvaluator{}}
	e.flowEvaluator = e
	return e
}

func (e *CompletedEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experiment := request.Experiment
	if experiment.Status == model.ExperimentStatusCompleted {
		winnerVariation, ok := experiment.WinnerVariation()
		if !ok {
			return nil, false, fmt.Errorf("winner variation [%d]", experiment.ID)
		}
		return e.evaluation(request, context, winnerVariation, decision.ReasonExperimentCompleted)
	} else {
		return nextFlow.Evaluate(request, context)
	}
}

type TargetEvaluator struct {
	*baseFlowEvaluator
	determiner TargetDeterminer
}

func NewTargetEvaluator(determiner TargetDeterminer) *TargetEvaluator {
	e := &TargetEvaluator{&baseFlowEvaluator{}, determiner}
	e.flowEvaluator = e
	return e
}

func (e *TargetEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	if request.Experiment.Type != model.ExperimentTypeAbTest {
		return nil, false, fmt.Errorf("experiment type must be AB_TEST [%d]", request.Experiment.ID)
	}
	isUserInExperimentTarget, err := e.determiner.IsUserInExperimentTarget(request, context)
	if err != nil {
		return nil, false, err
	}
	if isUserInExperimentTarget {
		return nextFlow.Evaluate(request, context)
	} else {
		return e.defaultEvaluation(request, context, decision.ReasonNotInExperimentTarget)
	}
}

type TrafficAllocateEvaluator struct {
	*baseFlowEvaluator
	actionResolver ActionResolver
}

func NewTrafficAllocatedEvaluator(actionResolver ActionResolver) *TrafficAllocateEvaluator {
	e := &TrafficAllocateEvaluator{&baseFlowEvaluator{}, actionResolver}
	e.flowEvaluator = e
	return e
}

func (e *TrafficAllocateEvaluator) evaluate(request Request, context evaluator.Context, _ flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experiment := request.Experiment
	if experiment.Status != model.ExperimentStatusRunning {
		return nil, false, fmt.Errorf("experiment status must be RUNNING [%d]", experiment.ID)
	}
	if experiment.Type != model.ExperimentTypeAbTest {
		return nil, false, fmt.Errorf("experiment type must be AB_TEST [%d]", experiment.ID)
	}
	defaultRule := experiment.DefaultRule
	variation, ok, err := e.actionResolver.Resolve(request, defaultRule)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return e.defaultEvaluation(request, context, decision.ReasonTrafficNotAllocated)
	}
	if variation.IsDropped {
		return e.defaultEvaluation(request, context, decision.ReasonVariationDropped)
	}
	return e.evaluation(request, context, variation, decision.ReasonTrafficAllocated)
}

type TargetRuleEvaluator struct {
	*baseFlowEvaluator
	determiner     TargetRuleDeterminer
	actionResolver ActionResolver
}

func NewTargetRuleEvaluator(determiner TargetRuleDeterminer, actionResolver ActionResolver) *TargetRuleEvaluator {
	e := &TargetRuleEvaluator{&baseFlowEvaluator{}, determiner, actionResolver}
	e.flowEvaluator = e
	return e
}

func (e *TargetRuleEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experiment := request.Experiment
	if experiment.Status != model.ExperimentStatusRunning {
		return nil, false, fmt.Errorf("experiment status must be RUNNING [%d]", experiment.ID)
	}
	if experiment.Type != model.ExperimentTypeFeatureFlag {
		return nil, false, fmt.Errorf("experiment type must be FEATURE_FLAG [%d]", experiment.ID)
	}
	if _, ok := request.user.Identifiers[experiment.IdentifierType]; !ok {
		return nextFlow.Evaluate(request, context)
	}
	targetRule, ok, err := e.determiner.Determine(request, context)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nextFlow.Evaluate(request, context)
	}
	variation, ok, err := e.actionResolver.Resolve(request, targetRule.Action)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, fmt.Errorf("feature flag must decide variation [%d]", experiment.ID)
	}
	return e.evaluation(request, context, variation, decision.ReasonTargetRuleMatch)
}

type DefaultRuleEvaluator struct {
	*baseFlowEvaluator
	actionResolver ActionResolver
}

func NewDefaultRuleEvaluator(actionResolver ActionResolver) *DefaultRuleEvaluator {
	e := &DefaultRuleEvaluator{&baseFlowEvaluator{}, actionResolver}
	e.flowEvaluator = e
	return e
}

func (e *DefaultRuleEvaluator) evaluate(request Request, context evaluator.Context, _ flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experiment := request.Experiment
	if experiment.Status != model.ExperimentStatusRunning {
		return nil, false, fmt.Errorf("experiment status must be RUNNING [%d]", experiment.ID)
	}
	if experiment.Type != model.ExperimentTypeFeatureFlag {
		return nil, false, fmt.Errorf("experiment type must be FEATURE_FLAG [%d]", experiment.ID)
	}
	if _, ok := request.user.Identifiers[experiment.IdentifierType]; !ok {
		return e.defaultEvaluation(request, context, decision.ReasonDefaultRule)
	}
	variation, ok, err := e.actionResolver.Resolve(request, experiment.DefaultRule)
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, fmt.Errorf("feature flag must decide variation [%d]", experiment.ID)
	}
	return e.evaluation(request, context, variation, decision.ReasonDefaultRule)
}

type ContainerEvaluator struct {
	*baseFlowEvaluator
	containerResolver ContainerResolver
}

func NewContainerEvaluator(containerResolver ContainerResolver) *ContainerEvaluator {
	e := &ContainerEvaluator{&baseFlowEvaluator{}, containerResolver}
	e.flowEvaluator = e
	return e
}

func (e *ContainerEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	experiment := request.Experiment
	if experiment.ContainerID == nil {
		return nextFlow.Evaluate(request, context)
	}

	container, ok := request.Workspace().GetContainer(*experiment.ContainerID)
	if !ok {
		return nil, false, fmt.Errorf("container [%d]", *experiment.ContainerID)
	}

	isUserInContainerGroup, err := e.containerResolver.IsUserInContainerGroup(request, container)
	if err != nil {
		return nil, false, err
	}
	if isUserInContainerGroup {
		return nextFlow.Evaluate(request, context)
	} else {
		return e.defaultEvaluation(request, context, decision.ReasonNotInMutualExclusionExperiment)
	}
}

type IdentifierEvaluator struct {
	*baseFlowEvaluator
}

func NewIdentifierEvaluator() *IdentifierEvaluator {
	e := &IdentifierEvaluator{&baseFlowEvaluator{}}
	e.flowEvaluator = e
	return e
}

func (e *IdentifierEvaluator) evaluate(request Request, context evaluator.Context, nextFlow flow.EvaluationFlow) (evaluator.Evaluation, bool, error) {
	if _, ok := request.user.Identifiers[request.Experiment.IdentifierType]; ok {
		return nextFlow.Evaluate(request, context)
	} else {
		return e.defaultEvaluation(request, context, decision.ReasonIdentifierNotFound)
	}
}
