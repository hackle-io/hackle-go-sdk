package remoteconfig

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/bucketer"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/target"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type TargetRuleDeterminer interface {
	Determine(request Request, context evaluator.Context) (model.RemoteConfigTargetRule, bool, error)
}

func NewTargetRuleDeterminer(targetMatcher target.Matcher, bucketer bucketer.Bucketer) TargetRuleDeterminer {
	matcher := &matcher{
		targetMatcher: targetMatcher,
		bucketer:      bucketer,
	}
	return &targetRuleDeterminer{
		matcher: matcher,
	}
}

type targetRuleDeterminer struct {
	matcher Matcher
}

func (d *targetRuleDeterminer) Determine(request Request, context evaluator.Context) (model.RemoteConfigTargetRule, bool, error) {
	for _, targetRule := range request.Parameter.TargetRules {
		matches, err := d.matcher.Matches(request, context, targetRule)
		if err != nil {
			return model.RemoteConfigTargetRule{}, false, err
		}
		if matches {
			return targetRule, true, nil
		}
	}
	return model.RemoteConfigTargetRule{}, false, nil
}

type Matcher interface {
	Matches(request Request, context evaluator.Context, targetRule model.RemoteConfigTargetRule) (bool, error)
}

type matcher struct {
	targetMatcher target.Matcher
	bucketer      bucketer.Bucketer
}

func (m *matcher) Matches(request Request, context evaluator.Context, targetRule model.RemoteConfigTargetRule) (bool, error) {
	matches, err := m.targetMatcher.Matches(request, context, targetRule.Target)
	if err != nil {
		return false, err
	}
	if !matches {
		return false, nil
	}
	identifier, ok := request.User().Identifiers[request.Parameter.IdentifierType]
	if !ok {
		return false, nil
	}
	bucket, ok := request.Workspace().GetBucket(targetRule.BucketID)
	if !ok {
		return false, fmt.Errorf("bucket [%d]", targetRule.BucketID)
	}
	_, ok = m.bucketer.Bucketing(bucket, identifier)
	return ok, nil
}
