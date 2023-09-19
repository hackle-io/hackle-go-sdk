package experiment

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/bucketer"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/target"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type ActionResolver interface {
	Resolve(request Request, action model.Action) (model.Variation, bool, error)
}

type actionResolver struct {
	bucketer bucketer.Bucketer
}

func (r *actionResolver) Resolve(request Request, action model.Action) (model.Variation, bool, error) {
	switch action.Type {
	case model.ActionTypeVariation:
		return r.resolveVariation(request, action)
	case model.ActionTypeBucket:
		return r.resolveBucket(request, action)
	}
	return model.Variation{}, false, fmt.Errorf("unsupported action type [%s]", action.Type)
}

func (r *actionResolver) resolveVariation(request Request, action model.Action) (model.Variation, bool, error) {
	variationID := action.VariationID
	if variationID == nil {
		return model.Variation{}, false, fmt.Errorf("action variation [%d]", request.Experiment.ID)
	}
	variation, ok := request.Experiment.GetVariationByID(*variationID)
	if !ok {
		return model.Variation{}, false, fmt.Errorf("variation [%d]", *variationID)
	}
	return variation, true, nil
}

func (r *actionResolver) resolveBucket(request Request, action model.Action) (model.Variation, bool, error) {
	bucketID := action.BucketID
	if bucketID == nil {
		return model.Variation{}, false, fmt.Errorf("action bucket [%d]", request.Experiment.ID)
	}
	bucket, ok := request.Workspace().GetBucket(*bucketID)
	if !ok {
		return model.Variation{}, false, fmt.Errorf("bucket [%d]", *bucketID)
	}
	identifier, ok := request.user.Identifiers[request.Experiment.IdentifierType]
	if !ok {
		return model.Variation{}, false, nil
	}
	slot, ok := r.bucketer.Bucketing(bucket, identifier)
	if !ok {
		return model.Variation{}, false, nil
	}
	variation, ok := request.Experiment.GetVariationByID(slot.VariationID)
	return variation, ok, nil
}

type OverrideResolver interface {
	Resolve(request Request, context evaluator.Context) (model.Variation, bool, error)
}

type overrideResolver struct {
	targetMatcher  target.Matcher
	actionResolver ActionResolver
}

func (r *overrideResolver) Resolve(request Request, context evaluator.Context) (model.Variation, bool, error) {
	if userOverriddenVariation, ok := r.resolveUserOverride(request); ok {
		return userOverriddenVariation, true, nil
	}
	return r.resolveSegmentOverride(request, context)
}

func (r *overrideResolver) resolveUserOverride(request Request) (model.Variation, bool) {
	experiment := request.Experiment
	identifier, ok := request.user.Identifiers[experiment.IdentifierType]
	if !ok {
		return model.Variation{}, false
	}
	overriddenVariationId, ok := experiment.UserOverrides[identifier]
	if !ok {
		return model.Variation{}, false
	}
	return experiment.GetVariationByID(overriddenVariationId)
}

func (r *overrideResolver) resolveSegmentOverride(request Request, context evaluator.Context) (model.Variation, bool, error) {
	for _, overriddenRule := range request.Experiment.SegmentOverrides {
		matches, err := r.targetMatcher.Matches(request, context, overriddenRule.Target)
		if err != nil {
			return model.Variation{}, false, err
		}
		if matches {
			return r.actionResolver.Resolve(request, overriddenRule.Action)
		}

	}
	return model.Variation{}, false, nil
}

type ContainerResolver interface {
	IsUserInContainerGroup(request Request, container model.Container) (bool, error)
}

type containerResolver struct {
	bucketer bucketer.Bucketer
}

func (r *containerResolver) IsUserInContainerGroup(request Request, container model.Container) (bool, error) {

	experiment := request.Experiment

	identifier, ok := request.user.Identifiers[experiment.IdentifierType]
	if !ok {
		return false, nil
	}

	bucket, ok := request.Workspace().GetBucket(container.BucketID)
	if !ok {
		return false, fmt.Errorf("bucket [%d]", container.BucketID)
	}

	slot, ok := r.bucketer.Bucketing(bucket, identifier)
	if !ok {
		return false, nil
	}

	group, ok := container.GetGroup(slot.VariationID)
	if !ok {
		return false, fmt.Errorf("container group [%d]", slot.VariationID)
	}

	for _, experimentId := range group.Experiments {
		if experimentId == experiment.ID {
			return true, nil
		}
	}
	return false, nil
}
