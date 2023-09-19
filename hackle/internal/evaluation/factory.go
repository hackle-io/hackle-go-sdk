package evaluation

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/bucketer"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/delegating"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/remoteconfig"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition"
	conditionexperiment "github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition/segment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/condition/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/target"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/value"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

func NewEvaluators() (experiment.Evaluator, remoteconfig.Evaluator) {

	delegatingEvaluator := delegating.NewEvaluator()

	targetMatcher := NewTargetMatcher(delegatingEvaluator)
	buckter := bucketer.NewBucketer()

	experimentEvaluator := experiment.NewEvaluator(experiment.NewFlowFactory(targetMatcher, buckter))
	delegatingEvaluator.Add(experimentEvaluator)

	remoteConfigEvaluator := remoteconfig.NewEvaluator(remoteconfig.NewTargetRuleDeterminer(targetMatcher, buckter))
	delegatingEvaluator.Add(remoteConfigEvaluator)

	return experimentEvaluator, remoteConfigEvaluator
}

func NewTargetMatcher(evaluator evaluator.Evaluator) target.Matcher {
	conditionMatcherFactory := NewConditionMatcherFactory(evaluator)
	return target.NewMatcher(conditionMatcherFactory)
}

func NewConditionMatcherFactory(evaluator evaluator.Evaluator) condition.MatcherFactory {
	valueOperatorMatcher := value.NewOperatorMatcher()
	userConditionMatcher := user.NewConditionMatcher(valueOperatorMatcher)
	segmentConditionMatcher := segment.NewConditionMatcher(userConditionMatcher)
	experimentConditionMatcher := conditionexperiment.NewConditionMatcher(evaluator, valueOperatorMatcher)
	return condition.NewMatcherFactory(map[model.TargetKeyType]condition.Matcher{
		model.TargetKeyTypeUserId:         userConditionMatcher,
		model.TargetKeyTypeUserProperty:   userConditionMatcher,
		model.TargetKeyTypeHackleProperty: userConditionMatcher,
		model.TargetKeyTypeSegment:        segmentConditionMatcher,
		model.TargetKeyTypeAbTest:         experimentConditionMatcher,
		model.TargetKeyTypeFeatureFlag:    experimentConditionMatcher,
	})
}
