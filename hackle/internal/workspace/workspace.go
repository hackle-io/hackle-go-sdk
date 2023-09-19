package workspace

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
)

type Workspace interface {
	GetExperiment(experimentKey int64) (model.Experiment, bool)
	GetFeatureFlag(featureKey int64) (model.Experiment, bool)
	GetEventType(eventTypeKey string) (model.EventType, bool)
	GetBucket(bucketID int64) (model.Bucket, bool)
	GetSegment(segmentKey string) (model.Segment, bool)
	GetContainer(containerID int64) (model.Container, bool)
	GetParameterConfiguration(parameterConfigurationID int64) (model.ParameterConfiguration, bool)
	GetRemoteConfigParameter(parameterKey string) (model.RemoteConfigParameter, bool)
}

type workspace struct {
	experiments             map[int64]model.Experiment
	featureFlags            map[int64]model.Experiment
	buckets                 map[int64]model.Bucket
	eventTypes              map[string]model.EventType
	segments                map[string]model.Segment
	containers              map[int64]model.Container
	parameterConfigurations map[int64]model.ParameterConfiguration
	remoteConfigParameters  map[string]model.RemoteConfigParameter
}

func NewFrom(dto WorkspaceDTO) Workspace {

	experiments := make([]model.Experiment, 0)
	for _, it := range dto.Experiments {
		if experiment, ok := newExperiment(it, model.ExperimentTypeAbTest); ok {
			experiments = append(experiments, experiment)
		}
	}

	featureFlags := make([]model.Experiment, 0)
	for _, it := range dto.FeatureFlags {
		if featureFlag, ok := newExperiment(it, model.ExperimentTypeFeatureFlag); ok {
			featureFlags = append(featureFlags, featureFlag)
		}
	}

	buckets := make([]model.Bucket, 0)
	for _, it := range dto.Buckets {
		buckets = append(buckets, newBucket(it))
	}

	eventTypes := make([]model.EventType, 0)
	for _, it := range dto.Events {
		eventTypes = append(eventTypes, newEventType(it))
	}

	segments := make([]model.Segment, 0)
	for _, it := range dto.Segments {
		if segment, ok := newSegment(it); ok {
			segments = append(segments, segment)
		}
	}

	containers := make([]model.Container, 0)
	for _, it := range dto.Containers {
		containers = append(containers, newContainer(it))
	}

	parameterConfigurations := make([]model.ParameterConfiguration, 0)
	for _, it := range dto.ParameterConfigurations {
		parameterConfigurations = append(parameterConfigurations, newParameterConfiguration(it))
	}

	remoteConfigParameters := make([]model.RemoteConfigParameter, 0)
	for _, it := range dto.RemoteConfigParameters {
		if remoteConfigParameter, ok := newRemoteConfigParameter(it); ok {
			remoteConfigParameters = append(remoteConfigParameters, remoteConfigParameter)
		}
	}

	return New(
		experiments,
		featureFlags,
		buckets,
		eventTypes,
		segments,
		containers,
		parameterConfigurations,
		remoteConfigParameters,
	)
}

func New(
	experiments []model.Experiment,
	featureFlags []model.Experiment,
	buckets []model.Bucket,
	eventTypes []model.EventType,
	segments []model.Segment,
	containers []model.Container,
	parameterConfigurations []model.ParameterConfiguration,
	remoteConfigParameters []model.RemoteConfigParameter,
) Workspace {

	_experiments := make(map[int64]model.Experiment)
	for _, it := range experiments {
		_experiments[it.Key] = it
	}

	_featureFlags := make(map[int64]model.Experiment)
	for _, it := range featureFlags {
		_featureFlags[it.Key] = it
	}

	_buckets := make(map[int64]model.Bucket)
	for _, it := range buckets {
		_buckets[it.ID] = it
	}

	_eventTypes := make(map[string]model.EventType)
	for _, it := range eventTypes {
		_eventTypes[it.Key] = it
	}

	_segments := make(map[string]model.Segment)
	for _, it := range segments {
		_segments[it.Key] = it
	}

	_containers := make(map[int64]model.Container)
	for _, it := range containers {
		_containers[it.ID] = it
	}

	_parameterConfigurations := make(map[int64]model.ParameterConfiguration)
	for _, it := range parameterConfigurations {
		_parameterConfigurations[it.ID] = it
	}

	_remoteConfigParameters := make(map[string]model.RemoteConfigParameter)
	for _, it := range remoteConfigParameters {
		_remoteConfigParameters[it.Key] = it
	}

	return &workspace{
		experiments:             _experiments,
		featureFlags:            _featureFlags,
		buckets:                 _buckets,
		eventTypes:              _eventTypes,
		segments:                _segments,
		containers:              _containers,
		parameterConfigurations: _parameterConfigurations,
		remoteConfigParameters:  _remoteConfigParameters,
	}
}

func (w *workspace) GetExperiment(experimentKey int64) (model.Experiment, bool) {
	experiment, ok := w.experiments[experimentKey]
	return experiment, ok
}

func (w *workspace) GetFeatureFlag(featureKey int64) (model.Experiment, bool) {
	featureFlag, ok := w.featureFlags[featureKey]
	return featureFlag, ok
}

func (w *workspace) GetEventType(eventTypeKey string) (model.EventType, bool) {
	eventType, ok := w.eventTypes[eventTypeKey]
	return eventType, ok
}

func (w *workspace) GetBucket(bucketID int64) (model.Bucket, bool) {
	bucket, ok := w.buckets[bucketID]
	return bucket, ok
}

func (w *workspace) GetSegment(segmentKey string) (model.Segment, bool) {
	segment, ok := w.segments[segmentKey]
	return segment, ok
}

func (w *workspace) GetContainer(containerID int64) (model.Container, bool) {
	container, ok := w.containers[containerID]
	return container, ok
}

func (w *workspace) GetParameterConfiguration(parameterConfigurationID int64) (model.ParameterConfiguration, bool) {
	parameterConfiguration, ok := w.parameterConfigurations[parameterConfigurationID]
	return parameterConfiguration, ok
}

func (w *workspace) GetRemoteConfigParameter(parameterKey string) (model.RemoteConfigParameter, bool) {
	remoteConfigParameter, ok := w.remoteConfigParameters[parameterKey]
	return remoteConfigParameter, ok
}

func newExperiment(dto ExperimentDTO, experimentType model.ExperimentType) (model.Experiment, bool) {
	execution := dto.Execution
	status, ok := model.NewExperimentStatusFrom(execution.Status)
	if !ok {
		return model.Experiment{}, false
	}

	variations := make([]model.Variation, 0)
	for _, it := range dto.Variations {
		variations = append(variations, newVariationFrom(it))
	}

	userOverrides := make(map[string]int64)
	for _, it := range execution.UserOverrides {
		userOverrides[it.UserID] = it.VariationID
	}

	segmentOverrides := make([]model.TargetRule, 0)
	for _, it := range execution.SegmentOverrides {
		if targetRule, ok := newTargetRule(it, model.TargetingTypeIdentifier); ok {
			segmentOverrides = append(segmentOverrides, targetRule)
		}
	}

	targetAudiences := make([]model.Target, 0)
	for _, it := range execution.TargetAudiences {
		if target, ok := newTarget(it, model.TargetingTypeProperty); ok {
			targetAudiences = append(targetAudiences, target)
		}
	}

	targetRules := make([]model.TargetRule, 0)
	for _, it := range execution.TargetRules {
		if targetRule, ok := newTargetRule(it, model.TargetingTypeProperty); ok {
			targetRules = append(targetRules, targetRule)
		}
	}

	defaultRule, ok := newTargetAction(execution.DefaultRule)
	if !ok {
		return model.Experiment{}, false
	}

	return model.Experiment{
		ID:                dto.ID,
		Key:               dto.Key,
		Name:              dto.Name,
		Type:              experimentType,
		IdentifierType:    dto.IdentifierType,
		Status:            status,
		Version:           dto.Version,
		ExecutionVersion:  dto.Execution.Version,
		Variations:        variations,
		UserOverrides:     userOverrides,
		SegmentOverrides:  segmentOverrides,
		TargetAudiences:   targetAudiences,
		TargetRules:       targetRules,
		DefaultRule:       defaultRule,
		ContainerID:       dto.ContainerID,
		WinnerVariationID: dto.WinnerVariationID,
	}, true
}

func newVariationFrom(dto VariationDTO) model.Variation {
	return model.Variation{
		ID:                       dto.ID,
		Key:                      dto.Key,
		IsDropped:                dto.Status == "DROPPED",
		ParameterConfigurationID: dto.ParameterConfigurationID,
	}
}

func newTarget(dto TargetDTO, targetingType model.TargetingType) (model.Target, bool) {
	conditions := make([]model.TargetCondition, 0)
	for _, it := range dto.Conditions {
		if condition, ok := newCondition(it, targetingType); ok {
			conditions = append(conditions, condition)
		}
	}
	if len(conditions) == 0 {
		return model.Target{}, false
	}
	return model.Target{
		Conditions: conditions,
	}, true
}

func newCondition(dto TargetConditionDTO, targetingType model.TargetingType) (model.TargetCondition, bool) {
	key, ok := newTargetKey(dto.Key)
	if !ok {
		return model.TargetCondition{}, false
	}
	if !targetingType.Supports(key.Type) {
		return model.TargetCondition{}, false
	}
	match, ok := newTargetMatch(dto.Match)
	if !ok {
		return model.TargetCondition{}, false
	}
	return model.TargetCondition{
		Key:   key,
		Match: match,
	}, true
}

func newTargetKey(dto TargetKeyDTO) (model.TargetKey, bool) {
	keyType, ok := model.TargetKeyTypeFrom(dto.Type)
	if !ok {
		return model.TargetKey{}, false
	}
	return model.TargetKey{
		Type: keyType,
		Name: dto.Name,
	}, true
}

func newTargetMatch(dto TargetMatchDTO) (model.TargetMatch, bool) {
	matchType, ok := model.TargetMatchTypeFrom(dto.Type)
	if !ok {
		return model.TargetMatch{}, false
	}
	operator, ok := model.TargetOperatorFrom(dto.Operator)
	if !ok {
		return model.TargetMatch{}, false
	}
	valueType, ok := types.TypeFrom(dto.ValueType)
	if !ok {
		return model.TargetMatch{}, false
	}
	return model.TargetMatch{
		Type:      matchType,
		Operator:  operator,
		ValueType: valueType,
		Values:    dto.Values,
	}, true
}

func newTargetAction(dto TargetActionDTO) (model.Action, bool) {
	actionType, ok := model.ActionTypeFrom(dto.Type)
	if !ok {
		return model.Action{}, false
	}
	return model.Action{
		Type:        actionType,
		VariationID: dto.VariationID,
		BucketID:    dto.BucketID,
	}, true
}

func newTargetRule(dto TargetRuleDTO, targetingType model.TargetingType) (model.TargetRule, bool) {
	target, ok := newTarget(dto.Target, targetingType)
	if !ok {
		return model.TargetRule{}, false
	}
	action, ok := newTargetAction(dto.Action)
	if !ok {
		return model.TargetRule{}, false
	}
	return model.TargetRule{
		Target: target,
		Action: action,
	}, true
}

func newBucket(dto BucketDTO) model.Bucket {
	slots := make([]model.Slot, 0)
	for _, it := range dto.Slots {
		slots = append(slots, newSlot(it))
	}
	return model.Bucket{
		ID:       dto.ID,
		Seed:     dto.Seed,
		SlotSize: dto.SlotSize,
		Slots:    slots,
	}
}

func newSlot(dto SlotDTO) model.Slot {
	return model.Slot{
		StartInclusive: dto.StartInclusive,
		EndExclusive:   dto.EndExclusive,
		VariationID:    dto.VariationID,
	}
}

func newEventType(dto EventTypeDTO) model.EventType {
	return model.EventType{
		ID:  dto.ID,
		Key: dto.Key,
	}
}

func newSegment(dto SegmentDTO) (model.Segment, bool) {
	segmentType, ok := model.SegmentTypeFrom(dto.Type)
	if !ok {
		return model.Segment{}, false
	}
	targets := make([]model.Target, 0)
	for _, it := range dto.Targets {
		if target, ok := newTarget(it, model.TargetingTypeSegment); ok {
			targets = append(targets, target)
		}
	}
	return model.Segment{
		ID:      dto.ID,
		Key:     dto.Key,
		Type:    segmentType,
		Targets: targets,
	}, true
}

func newContainer(dto ContainerDTO) model.Container {
	groups := make([]model.ContainerGroup, 0)
	for _, it := range dto.Groups {
		groups = append(groups, newContainerGroup(it))
	}
	return model.Container{
		ID:       dto.ID,
		BucketID: dto.BucketID,
		Groups:   groups,
	}
}

func newContainerGroup(dto ContainerGroupDTO) model.ContainerGroup {
	return model.ContainerGroup{
		ID:          dto.ID,
		Experiments: dto.Experiments,
	}
}

func newParameterConfiguration(dto ParameterConfigurationDTO) model.ParameterConfiguration {
	parameters := make(map[string]interface{})
	for _, it := range dto.Parameters {
		parameters[it.Key] = it.Value
	}
	return model.ParameterConfiguration{
		ID:         dto.ID,
		Parameters: parameters,
	}
}

func newRemoteConfigParameter(dto RemoteConfigParameterDTO) (model.RemoteConfigParameter, bool) {
	valueType, ok := types.TypeFrom(dto.Type)
	if !ok {
		return model.RemoteConfigParameter{}, false
	}

	targetRules := make([]model.RemoteConfigTargetRule, 0)
	for _, it := range dto.TargetRules {
		if targetRule, ok := newRemoteConfigTargetRule(it); ok {
			targetRules = append(targetRules, targetRule)
		}
	}

	return model.RemoteConfigParameter{
		ID:             dto.ID,
		Key:            dto.Key,
		Type:           valueType,
		IdentifierType: dto.IdentifierType,
		TargetRules:    targetRules,
		DefaultValue: model.RemoteConfigValue{
			ID:       dto.DefaultValue.ID,
			RawValue: dto.DefaultValue.Value,
		},
	}, true
}

func newRemoteConfigTargetRule(dto RemoteConfigTargetRuleDTO) (model.RemoteConfigTargetRule, bool) {
	target, ok := newTarget(dto.Target, model.TargetingTypeProperty)
	if !ok {
		return model.RemoteConfigTargetRule{}, false
	}
	return model.RemoteConfigTargetRule{
		Key:      dto.Key,
		Name:     dto.Name,
		Target:   target,
		BucketID: dto.BucketID,
		Value: model.RemoteConfigValue{
			ID:       dto.Value.ID,
			RawValue: dto.Value.Value,
		},
	}, true
}
