package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/value"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"strconv"
)

type experimentMatcher interface {
	matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error)
}

type evaluatorMatcher interface {
	experimentMatcher
	getExperiment(request evaluator.Request, key int64) (model.Experiment, bool)
	resolveEvaluation(request evaluator.Request, evaluation experiment.Evaluation) experiment.Evaluation
	evaluationMatches(evaluation experiment.Evaluation, condition model.TargetCondition) bool
}

type baseEvaluatorMatcher struct {
	evaluatorMatcher
	evaluator evaluator.Evaluator
	matcher   value.OperatorMatcher
}

func (m *baseEvaluatorMatcher) matches(request evaluator.Request, context evaluator.Context, condition model.TargetCondition) (bool, error) {
	key, err := strconv.ParseInt(condition.Key.Name, 10, 64)
	if err != nil {
		return false, fmt.Errorf("invalid key [%s, %s]", condition.Key.Type, condition.Key.Name)
	}

	exp, ok := m.getExperiment(request, key)
	if !ok {
		return false, nil
	}

	evaluation, ok := getEvaluation(context, exp)
	if !ok {
		evaluation, err = m.evaluate(request, context, exp)
		if err != nil {
			return false, err
		}
	}
	experimentEvaluation, ok := evaluation.(experiment.Evaluation)

	if !ok {
		return false, fmt.Errorf("unexpected evaluation: %T (expected: experiment.Evaluation)", evaluation)
	}

	return m.evaluationMatches(experimentEvaluation, condition), nil
}

func (m *baseEvaluatorMatcher) evaluate(
	request evaluator.Request,
	context evaluator.Context,
	exp model.Experiment,
) (evaluator.Evaluation, error) {
	experimentRequest := experiment.NewRequestFrom(request, exp)
	evaluation, err := m.evaluator.Evaluate(experimentRequest, context)
	if err != nil {
		return nil, err
	}
	experimentEvaluation, ok := evaluation.(experiment.Evaluation)
	if !ok {
		return nil, fmt.Errorf("unexpected evaluation: %T (expected: experiment.Evaluation)", evaluation)
	}
	resolvedEvaluation := m.resolveEvaluation(request, experimentEvaluation)
	context.AddEvaluation(resolvedEvaluation)
	return resolvedEvaluation, nil
}

func getEvaluation(context evaluator.Context, exp model.Experiment) (evaluator.Evaluation, bool) {
	for _, evaluation := range context.Evaluations() {
		if experimentEvaluation, ok := evaluation.(experiment.Evaluation); ok && experimentEvaluation.Experiment.ID == exp.ID {
			return experimentEvaluation, true
		}
	}
	return nil, false
}

var abtestMatchedReasons = []string{
	decision.ReasonOverridden,
	decision.ReasonTrafficAllocated,
	decision.ReasonExperimentCompleted,
	decision.ReasonTrafficAllocatedByTargeting,
}

type abTestEvaluatorMatcher struct {
	*baseEvaluatorMatcher
}

func newAbTestEvaluatorMatcher(evaluator evaluator.Evaluator, matcher value.OperatorMatcher) *abTestEvaluatorMatcher {
	base := &baseEvaluatorMatcher{evaluator: evaluator, matcher: matcher}
	m := &abTestEvaluatorMatcher{base}
	m.evaluatorMatcher = m
	return m
}

func (m *abTestEvaluatorMatcher) getExperiment(request evaluator.Request, key int64) (model.Experiment, bool) {
	return request.Workspace().GetExperiment(key)
}

func (m *abTestEvaluatorMatcher) resolveEvaluation(request evaluator.Request, evaluation experiment.Evaluation) experiment.Evaluation {
	_, ok := request.(experiment.Request)
	if ok && evaluation.Reason() == decision.ReasonTrafficAllocated {
		return evaluation.With(decision.ReasonTrafficAllocatedByTargeting)
	}
	return evaluation
}

func (m *abTestEvaluatorMatcher) evaluationMatches(evaluation experiment.Evaluation, condition model.TargetCondition) bool {
	for _, reason := range abtestMatchedReasons {
		if reason == evaluation.Reason() {
			return m.matcher.Matches(evaluation.VariationKey, condition.Match)
		}
	}
	return false
}

type featureFlagEvaluatorMatcher struct {
	*baseEvaluatorMatcher
}

func newFeatureFlagEvaluatorMatcher(evaluator evaluator.Evaluator, matcher value.OperatorMatcher) *featureFlagEvaluatorMatcher {
	base := &baseEvaluatorMatcher{evaluator: evaluator, matcher: matcher}
	m := &featureFlagEvaluatorMatcher{base}
	m.evaluatorMatcher = m
	return m
}

func (m *featureFlagEvaluatorMatcher) getExperiment(request evaluator.Request, key int64) (model.Experiment, bool) {
	return request.Workspace().GetFeatureFlag(key)
}

func (m *featureFlagEvaluatorMatcher) resolveEvaluation(_ evaluator.Request, evaluation experiment.Evaluation) experiment.Evaluation {
	return evaluation
}

func (m *featureFlagEvaluatorMatcher) evaluationMatches(evaluation experiment.Evaluation, condition model.TargetCondition) bool {
	on := evaluation.VariationKey != "A"
	return m.matcher.Matches(on, condition.Match)
}
