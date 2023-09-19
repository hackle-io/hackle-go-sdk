package mocks

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
)

type MockWorkspace struct {
	experiments             map[int64]model.Experiment
	featureFlags            map[int64]model.Experiment
	buckets                 map[int64]model.Bucket
	eventTypes              map[string]model.EventType
	segments                map[string]model.Segment
	containers              map[int64]model.Container
	parameterConfigurations map[int64]model.ParameterConfiguration
	remoteConfigParameters  map[string]model.RemoteConfigParameter
}

func CreateWorkspace() *MockWorkspace {
	return &MockWorkspace{
		experiments:             make(map[int64]model.Experiment),
		featureFlags:            make(map[int64]model.Experiment),
		buckets:                 make(map[int64]model.Bucket),
		eventTypes:              make(map[string]model.EventType),
		segments:                make(map[string]model.Segment),
		containers:              make(map[int64]model.Container),
		parameterConfigurations: make(map[int64]model.ParameterConfiguration),
		remoteConfigParameters:  make(map[string]model.RemoteConfigParameter),
	}
}

func (m *MockWorkspace) GetExperiment(experimentKey int64) (model.Experiment, bool) {
	experiment, ok := m.experiments[experimentKey]
	return experiment, ok
}

func (m *MockWorkspace) GetFeatureFlag(featureKey int64) (model.Experiment, bool) {
	featureFlag, ok := m.featureFlags[featureKey]
	return featureFlag, ok
}

func (m *MockWorkspace) GetEventType(eventTypeKey string) (model.EventType, bool) {
	eventType, ok := m.eventTypes[eventTypeKey]
	return eventType, ok
}

func (m *MockWorkspace) GetBucket(bucketID int64) (model.Bucket, bool) {
	bucket, ok := m.buckets[bucketID]
	return bucket, ok
}

func (m *MockWorkspace) GetSegment(segmentKey string) (model.Segment, bool) {
	segment, ok := m.segments[segmentKey]
	return segment, ok
}

func (m *MockWorkspace) GetContainer(containerID int64) (model.Container, bool) {
	container, ok := m.containers[containerID]
	return container, ok
}

func (m *MockWorkspace) GetParameterConfiguration(parameterConfigurationID int64) (model.ParameterConfiguration, bool) {
	parameterConfiguration, ok := m.parameterConfigurations[parameterConfigurationID]
	return parameterConfiguration, ok
}

func (m *MockWorkspace) GetRemoteConfigParameter(parameterKey string) (model.RemoteConfigParameter, bool) {
	remoteConfigParameter, ok := m.remoteConfigParameters[parameterKey]
	return remoteConfigParameter, ok
}

func (m *MockWorkspace) Experiment(experiment model.Experiment) *MockWorkspace {
	m.experiments[experiment.Key] = experiment
	return m
}

func (m *MockWorkspace) FeatureFlag(featureFlag model.Experiment) *MockWorkspace {
	m.featureFlags[featureFlag.Key] = featureFlag
	return m
}

func (m *MockWorkspace) Bucket(bucket model.Bucket) *MockWorkspace {
	m.buckets[bucket.ID] = bucket
	return m
}

func (m *MockWorkspace) RemoteConfigParameter(parameter model.RemoteConfigParameter) *MockWorkspace {
	m.remoteConfigParameters[parameter.Key] = parameter
	return m
}

func (m *MockWorkspace) EventType(eventType model.EventType) *MockWorkspace {
	m.eventTypes[eventType.Key] = eventType
	return m
}

func (m *MockWorkspace) Container(container model.Container) *MockWorkspace {
	m.containers[container.ID] = container
	return m
}

func (m *MockWorkspace) Segment(segment model.Segment) *MockWorkspace {
	m.segments[segment.Key] = segment
	return m
}

func (m *MockWorkspace) ParameterConfiguration(config model.ParameterConfiguration) *MockWorkspace {
	m.parameterConfigurations[config.ID] = config
	return m
}
