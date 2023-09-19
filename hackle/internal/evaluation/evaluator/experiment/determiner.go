package experiment

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/target"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type TargetDeterminer interface {
	IsUserInExperimentTarget(request Request, context evaluator.Context) (bool, error)
}

type targetDeterminer struct {
	targetMatcher target.Matcher
}

func (d *targetDeterminer) IsUserInExperimentTarget(request Request, context evaluator.Context) (bool, error) {
	targetAudiences := request.Experiment.TargetAudiences
	if len(targetAudiences) == 0 {
		return true, nil
	}
	for _, audience := range targetAudiences {
		matches, err := d.targetMatcher.Matches(request, context, audience)
		if err != nil {
			return false, err
		}
		if matches {
			return matches, nil
		}
	}
	return false, nil
}

type TargetRuleDeterminer interface {
	Determine(request Request, context evaluator.Context) (model.TargetRule, bool, error)
}

type targetRuleDeterminer struct {
	targetMatcher target.Matcher
}

func (d *targetRuleDeterminer) Determine(request Request, context evaluator.Context) (model.TargetRule, bool, error) {
	for _, targetRule := range request.Experiment.TargetRules {
		matches, err := d.targetMatcher.Matches(request, context, targetRule.Target)
		if err != nil {
			return model.TargetRule{}, false, err
		}
		if matches {
			return targetRule, true, nil
		}
	}
	return model.TargetRule{}, false, nil
}
