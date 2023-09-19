package core

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/remoteconfig"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/event"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

type fields struct {
	experimentEvaluator   *mockExperimentEvaluator
	remoteConfigEvaluator *mockRemoteConfigEvaluator
	workspaceFetcher      *mockWorkspaceFetcher
	eventFactory          *mockEventFactory
	eventProcessor        *mockEventProcessor
	clock                 clock.Clock
}

func sut() (*core, *fields) {
	eventProcessor := &mockEventProcessor{}
	eventProcessor.On("Process", mock.Anything).Return()
	f := &fields{
		experimentEvaluator:   &mockExperimentEvaluator{},
		remoteConfigEvaluator: &mockRemoteConfigEvaluator{},
		workspaceFetcher:      &mockWorkspaceFetcher{},
		eventFactory:          &mockEventFactory{},
		eventProcessor:        eventProcessor,
		clock:                 clock.Fixed(42),
	}
	return &core{
		experimentEvaluator:   f.experimentEvaluator,
		remoteConfigEvaluator: f.remoteConfigEvaluator,
		workspaceFetcher:      f.workspaceFetcher,
		eventFactory:          f.eventFactory,
		eventProcessor:        f.eventProcessor,
		clock:                 f.clock,
	}, f
}

func TestCore_Experiment(t *testing.T) {

	t.Run("when sdk not ready then return default variation", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(nil, false)

		// when
		actual, err := sut.Experiment(42, user.HackleUser{}, "A")

		// then
		assert.Nil(t, err)
		assert.Equal(t, "A", actual.Variation())
		assert.Equal(t, "SDK_NOT_READY", actual.Reason())
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when experiment not found then return default variation", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(mocks.CreateWorkspace(), false)

		// when
		actual, err := sut.Experiment(42, user.HackleUser{}, "A")

		// then
		assert.Nil(t, err)
		assert.Equal(t, "A", actual.Variation())
		assert.Equal(t, "EXPERIMENT_NOT_FOUND", actual.Reason())
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when error on evaluate then return error", func(t *testing.T) {
		// given
		sut, f := sut()
		ws := mocks.CreateWorkspace()
		ws.Experiment(model.Experiment{Key: 42})
		f.workspaceFetcher.On("Fetch").Return(ws, true)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(experiment.Evaluation{}, errors.New("fail"))

		// when
		_, actual := sut.Experiment(42, user.HackleUser{}, "A")

		// then
		assert.Equal(t, errors.New("fail"), actual)
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when error on create events then return error", func(t *testing.T) {
		// given
		sut, f := sut()
		exp := model.Experiment{Key: 42}
		ws := mocks.CreateWorkspace()
		ws.Experiment(exp)
		f.workspaceFetcher.On("Fetch").Return(ws, true)

		cfg := model.ParameterConfiguration{ID: 420, Parameters: map[string]interface{}{"A": "B"}}
		eval := experiment.NewEvaluationOf(
			decision.ReasonTrafficAllocated,
			make([]evaluator.Evaluation, 0),
			exp,
			ref.Int64(320),
			"B",
			&cfg,
		)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.On("Create", mock.Anything, mock.Anything).Return([]event.UserEvent{}, errors.New("fail"))

		// when
		_, err := sut.Experiment(42, user.HackleUser{}, "A")

		// then
		assert.Equal(t, errors.New("fail"), err)
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when experiment evaluated then return evaluated variation with process events", func(t *testing.T) {
		sut, f := sut()
		exp := model.Experiment{Key: 42}
		ws := mocks.CreateWorkspace()
		ws.Experiment(exp)
		f.workspaceFetcher.On("Fetch").Return(ws, true)

		eval := experiment.NewEvaluationOf(
			decision.ReasonTrafficAllocated,
			make([]evaluator.Evaluation, 0),
			exp,
			ref.Int64(320),
			"B",
			&model.ParameterConfiguration{ID: 420, Parameters: map[string]interface{}{"A": "B"}},
		)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.On("Create", mock.Anything, mock.Anything).Return([]event.UserEvent{event.ExposureEvent{}, event.ExposureEvent{}}, nil)
		f.eventProcessor.On("Process", mock.Anything).Return()

		// when
		actual, _ := sut.Experiment(42, user.HackleUser{}, "A")

		// then
		assert.Equal(t, "B", actual.Variation())
		assert.Equal(t, "TRAFFIC_ALLOCATED", actual.Reason())
		assert.Equal(t, "B", actual.GetString("A", "default"))
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 2)
	})
}

func TestCore_FeatureFlag(t *testing.T) {
	t.Run("when sdk not ready then return false", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(nil, false)

		// when
		actual, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, false, actual.IsOn())
		assert.Equal(t, "SDK_NOT_READY", actual.Reason())
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when feature flag not found then return false", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(mocks.CreateWorkspace(), false)

		// when
		actual, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, false, actual.IsOn())
		assert.Equal(t, "FEATURE_FLAG_NOT_FOUND", actual.Reason())
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when error on evaluate then return error", func(t *testing.T) {
		// given
		sut, f := sut()
		ws := mocks.CreateWorkspace()
		ws.FeatureFlag(model.Experiment{Key: 42})
		f.workspaceFetcher.On("Fetch").Return(ws, true)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(experiment.Evaluation{}, errors.New("fail"))

		// when
		_, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Equal(t, errors.New("fail"), err)
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 0)
	})

	t.Run("when error on create event then return error", func(t *testing.T) {
		// given
		sut, f := sut()
		ws := mocks.CreateWorkspace()
		flag := model.Experiment{Key: 42}
		ws.FeatureFlag(flag)
		f.workspaceFetcher.On("Fetch").Return(ws, true)
		eval := experiment.NewEvaluationOf(decision.ReasonTargetRuleMatch, make([]evaluator.Evaluation, 0), flag, ref.Int64(320), "B", nil)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.MockCreateError(errors.New("failed to create events"))

		// when
		_, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Equal(t, errors.New("failed to create events"), err)
	})

	t.Run("when feature flag evaluated as A then return false", func(t *testing.T) {
		// given
		sut, f := sut()
		ws := mocks.CreateWorkspace()
		flag := model.Experiment{Key: 42}
		ws.FeatureFlag(flag)
		f.workspaceFetcher.On("Fetch").Return(ws, true)
		eval := experiment.NewEvaluationOf(decision.ReasonDefaultRule, make([]evaluator.Evaluation, 0), flag, ref.Int64(320), "A", &model.ParameterConfiguration{ID: 420, Parameters: map[string]interface{}{"A": "B"}})
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.MockCreateReturn([]event.UserEvent{event.ExposureEvent{}, event.ExposureEvent{}})

		// when
		actual, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, false, actual.IsOn())
		assert.Equal(t, "DEFAULT_RULE", actual.Reason())
		assert.Equal(t, "B", actual.GetString("A", "default"))
	})

	t.Run("when feature flag evaluated as not A then return true", func(t *testing.T) {
		// given
		sut, f := sut()
		ws := mocks.CreateWorkspace()
		flag := model.Experiment{Key: 42}
		ws.FeatureFlag(flag)
		f.workspaceFetcher.On("Fetch").Return(ws, true)
		eval := experiment.NewEvaluationOf(decision.ReasonTargetRuleMatch, make([]evaluator.Evaluation, 0), flag, ref.Int64(320), "B", nil)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.MockCreateReturn([]event.UserEvent{event.ExposureEvent{}, event.ExposureEvent{}})

		// when
		actual, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Nil(t, err)
		assert.Equal(t, true, actual.IsOn())
		assert.Equal(t, "TARGET_RULE_MATCH", actual.Reason())
	})

	t.Run("when feature flag evaluated then send exposure events", func(t *testing.T) {
		// given
		sut, f := sut()
		ws := mocks.CreateWorkspace()
		flag := model.Experiment{Key: 42}
		ws.FeatureFlag(flag)
		f.workspaceFetcher.On("Fetch").Return(ws, true)
		eval := experiment.NewEvaluationOf(decision.ReasonTargetRuleMatch, make([]evaluator.Evaluation, 0), flag, ref.Int64(320), "B", nil)
		f.experimentEvaluator.On("EvaluateExperiment", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.MockCreateReturn([]event.UserEvent{event.ExposureEvent{}, event.ExposureEvent{}})

		// when
		_, err := sut.FeatureFlag(42, user.HackleUser{})

		// then
		assert.Nil(t, err)
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 2)
	})
}

func TestCore_RemoteConfig(t *testing.T) {

	t.Run("when sdk not ready then return default value", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(nil, false)

		// when
		actual, err := sut.RemoteConfig("42", user.HackleUser{}, types.String, "default")

		// then
		assert.Nil(t, err)
		assert.Equal(t, "default", actual.Value())
		assert.Equal(t, "SDK_NOT_READY", actual.Reason())
	})

	t.Run("when rc parameter not found then return default value", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(mocks.CreateWorkspace(), true)

		// when
		actual, err := sut.RemoteConfig("42", user.HackleUser{}, types.String, "default")

		// then
		assert.Nil(t, err)
		assert.Equal(t, "default", actual.Value())
		assert.Equal(t, "REMOTE_CONFIG_PARAMETER_NOT_FOUND", actual.Reason())
	})

	t.Run("when error on evaluate rc then return error", func(t *testing.T) {
		// given
		sut, f := sut()
		parameter := model.RemoteConfigParameter{Key: "42"}
		ws := mocks.CreateWorkspace()
		ws.RemoteConfigParameter(parameter)
		f.workspaceFetcher.On("Fetch").Return(ws, true)

		f.remoteConfigEvaluator.On("EvaluateRemoteConfig", mock.Anything, mock.Anything).Return(remoteconfig.Evaluation{}, errors.New("failed to evaluate"))

		// when
		_, err := sut.RemoteConfig("42", user.HackleUser{}, types.String, "default")

		// then
		assert.Equal(t, errors.New("failed to evaluate"), err)
	})

	t.Run("when error on create events then return error", func(t *testing.T) {
		// given
		sut, f := sut()
		parameter := model.RemoteConfigParameter{Key: "42"}
		ws := mocks.CreateWorkspace()
		ws.RemoteConfigParameter(parameter)
		f.workspaceFetcher.On("Fetch").Return(ws, true)

		f.remoteConfigEvaluator.On("EvaluateRemoteConfig", mock.Anything, mock.Anything).Return(remoteconfig.Evaluation{}, nil)
		f.eventFactory.MockCreateError(errors.New("failed to create events"))

		// when
		_, err := sut.RemoteConfig("42", user.HackleUser{}, types.String, "default")

		// then
		assert.Equal(t, errors.New("failed to create events"), err)
	})

	t.Run("when rc evaluated then return evaluated value with send events", func(t *testing.T) {
		// given
		sut, f := sut()
		parameter := model.RemoteConfigParameter{Key: "42"}
		ws := mocks.CreateWorkspace()
		ws.RemoteConfigParameter(parameter)
		f.workspaceFetcher.On("Fetch").Return(ws, true)

		eval := remoteconfig.NewEvaluationOf(
			decision.ReasonDefaultRule,
			make([]evaluator.Evaluation, 0),
			parameter,
			ref.Int64(320),
			"evaluated",
			make(map[string]interface{}),
		)
		f.remoteConfigEvaluator.On("EvaluateRemoteConfig", mock.Anything, mock.Anything).Return(eval, nil)
		f.eventFactory.MockCreateReturn([]event.UserEvent{event.RemoteConfigEvent{}, event.ExposureEvent{}})

		// when
		actual, err := sut.RemoteConfig("42", user.HackleUser{}, types.String, "default")

		// then
		assert.Nil(t, err)
		assert.Equal(t, "evaluated", actual.Value())
		assert.Equal(t, "DEFAULT_RULE", actual.Reason())
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 2)
	})
}

func TestCore_Track(t *testing.T) {

	t.Run("when sdk not ready then track undefined event", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(nil, false)

		// when
		sut.Track(mocks.CreateEvent("42"), user.HackleUser{})

		// then
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 1)
		trackEvent := f.eventProcessor.Calls[0].Arguments[0].(event.TrackEvent)
		assert.Equal(t, int64(0), trackEvent.EventType.ID)
	})

	t.Run("when event type not found then track undefined event", func(t *testing.T) {
		// given
		sut, f := sut()
		f.workspaceFetcher.On("Fetch").Return(mocks.CreateWorkspace(), true)

		// when
		sut.Track(mocks.CreateEvent("42"), user.HackleUser{})

		// then
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 1)
		trackEvent := f.eventProcessor.Calls[0].Arguments[0].(event.TrackEvent)
		assert.Equal(t, int64(0), trackEvent.EventType.ID)
	})

	t.Run("track", func(t *testing.T) {
		// given
		sut, f := sut()
		eventType := model.EventType{ID: 42, Key: "42"}
		ws := mocks.CreateWorkspace()
		ws.EventType(eventType)
		f.workspaceFetcher.On("Fetch").Return(ws, true)

		// when
		sut.Track(mocks.CreateEvent("42"), user.HackleUser{})

		// then
		f.eventProcessor.AssertNumberOfCalls(t, "Process", 1)
		trackEvent := f.eventProcessor.Calls[0].Arguments[0].(event.TrackEvent)
		assert.Equal(t, int64(42), trackEvent.EventType.ID)
	})
}

func TestCore_Close(t *testing.T) {
	sut, f := sut()
	f.eventProcessor.On("Close").Return()
	f.workspaceFetcher.On("Close").Return()
	sut.Close()
	f.eventProcessor.AssertNumberOfCalls(t, "Close", 1)
	f.workspaceFetcher.AssertNumberOfCalls(t, "Close", 1)
}

func TestCore(t *testing.T) {

	/*
	 *       RC(1)
	 *      /     \
	 *     /       \
	 *  AB(2)     FF(4)
	 *    |   \     |
	 *    |     \   |
	 *  AB(3)     FF(5)
	 *              |
	 *              |
	 *            AB(6)
	 */
	t.Run("target_experiment", func(t *testing.T) {
		fetcher := workspace.NewFileFetcher("../../../testdata/workspace_target_experiment.json")
		processor := &memoryEventProcessor{}
		core := New(fetcher, processor)
		hackleUser := user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build()

		actual, err := core.RemoteConfig("rc", hackleUser, types.String, "!!")

		assert.Nil(t, err)
		assert.Equal(t, decision.NewRemoteConfigDecision("Targeting!!", decision.ReasonTargetRuleMatch), actual)

		assert.Equal(t, 6, len(processor.events))
		assert.Equal(t, map[string]interface{}{
			"requestValueType":    "STRING",
			"requestDefaultValue": "!!",
			"targetRuleKey":       "rc_1_key",
			"targetRuleName":      "rc_1_name",
			"returnValue":         "Targeting!!",
		}, (processor.events[0].(event.RemoteConfigEvent)).Properties)

		for _, userEvent := range processor.events[1:] {
			properties := (userEvent.(event.ExposureEvent)).Properties
			assert.Equal(t, "REMOTE_CONFIG", properties["$targetingRootType"])
			assert.Equal(t, int64(1), properties["$targetingRootId"])
		}
	})

	/*
	 *     RC(1)
	 *      ↓
	 * ┌── AB(2)
	 * ↑    ↓
	 * |   FF(3)
	 * ↑    ↓
	 * |   AB(4)
	 * └────┘
	 */
	t.Run("target_experiment_circular", func(t *testing.T) {
		fetcher := workspace.NewFileFetcher("../../../testdata/workspace_target_experiment_circular.json")
		processor := &memoryEventProcessor{}
		core := New(fetcher, processor)
		hackleUser := user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "a").Build()

		_, err := core.RemoteConfig("rc", hackleUser, types.String, "!!")

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "circular evaluation has occurred")
	})

	/*
	 *                     Container(1)
	 * ┌──────────────┬───────────────────────────────────────┐
	 * | ┌──────────┐ |                                       |
	 * | |   AB(2)  | |                                       |
	 * | └──────────┘ |                                       |
	 * └──────────────┴───────────────────────────────────────┘
	 *       25 %                        75 %
	 */
	t.Run("container", func(t *testing.T) {
		fetcher := workspace.NewFileFetcher("../../../testdata/workspace_container.json")
		processor := &memoryEventProcessor{}
		core := New(fetcher, processor)

		decisions := make([]decision.ExperimentDecision, 0)
		for i := 0; i < 10000; i++ {
			u := user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, strconv.Itoa(i)).Build()
			d, err := core.Experiment(2, u, "A")
			assert.Nil(t, err)
			decisions = append(decisions, d)
		}

		assert.Len(t, processor.events, 10000)
		assert.Len(t, decisions, 10000)

		count := func(reason string) int {
			var c int
			for _, d := range decisions {
				if d.Reason() == reason {
					c++
				}
			}
			return c
		}

		assert.Equal(t, 2452, count("TRAFFIC_ALLOCATED"))
		assert.Equal(t, 7548, count("NOT_IN_MUTUAL_EXCLUSION_EXPERIMENT"))
	})

	t.Run("segment_match", func(t *testing.T) {
		fetcher := workspace.NewFileFetcher("../../../testdata/workspace_segment_match.json")
		processor := &memoryEventProcessor{}
		core := New(fetcher, processor)

		d1, _ := core.Experiment(1, user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "matched_id").Build(), "A")
		assert.Equal(t, "A", d1.Variation())
		assert.Equal(t, "OVERRIDDEN", d1.Reason())

		d2, _ := core.Experiment(1, user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "not_matched_id").Build(), "A")
		assert.Equal(t, "A", d2.Variation())
		assert.Equal(t, "TRAFFIC_ALLOCATED", d2.Reason())
	})
}

type memoryEventProcessor struct {
	events []event.UserEvent
}

func (p *memoryEventProcessor) Process(event event.UserEvent) {
	p.events = append(p.events, event)
}

func (p *memoryEventProcessor) Start() {}

func (p *memoryEventProcessor) Close() {}

type mockExperimentEvaluator struct{ mock.Mock }

func (m *mockExperimentEvaluator) EvaluateExperiment(request experiment.Request, context evaluator.Context) (experiment.Evaluation, error) {
	arguments := m.Called(request, context)
	return arguments.Get(0).(experiment.Evaluation), arguments.Error(1)
}

type mockRemoteConfigEvaluator struct{ mock.Mock }

func (m *mockRemoteConfigEvaluator) EvaluateRemoteConfig(request remoteconfig.Request, context evaluator.Context) (remoteconfig.Evaluation, error) {
	arguments := m.MethodCalled("EvaluateRemoteConfig", request, context)
	return arguments.Get(0).(remoteconfig.Evaluation), arguments.Error(1)
}

type mockWorkspaceFetcher struct{ mock.Mock }

func (m *mockWorkspaceFetcher) Fetch() (workspace.Workspace, bool) {
	arguments := m.Called()
	if ws, ok := arguments.Get(0).(workspace.Workspace); ok {
		return ws, true
	}
	return nil, false
}

func (m *mockWorkspaceFetcher) Close() {
	m.Called()
}

type mockEventFactory struct{ mock.Mock }

func (m *mockEventFactory) Create(request evaluator.Request, evaluation evaluator.Evaluation) ([]event.UserEvent, error) {
	arguments := m.Called(request, evaluation)
	return arguments.Get(0).([]event.UserEvent), arguments.Error(1)
}

func (m *mockEventFactory) MockCreateReturn(events []event.UserEvent) {
	m.On("Create", mock.Anything, mock.Anything).Return(events, nil)
}

func (m *mockEventFactory) MockCreateError(err error) {
	m.On("Create", mock.Anything, mock.Anything).Return([]event.UserEvent{}, err)
}

type mockEventProcessor struct {
	mock.Mock
}

func (m *mockEventProcessor) Process(event event.UserEvent) {
	m.Called(event)
}

func (m *mockEventProcessor) Start() {
	m.Called()
}

func (m *mockEventProcessor) Close() {
	m.Called()
}
