package event

import (
	"github.com/google/uuid"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/experiment"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator/remoteconfig"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
)

type HackleEvent interface {
	Key() string
	Value() float64
	Properties() map[string]interface{}
}

type event struct {
	key        string
	value      float64
	properties map[string]interface{}
}

func (e event) Key() string {
	return e.key
}

func (e event) Value() float64 {
	return e.value
}

func (e event) Properties() map[string]interface{} {
	return e.properties
}

type UserEvent interface {
	InsertID() string
	Timestamp() int64
	User() user.HackleUser
}

func NewUserEvent(user user.HackleUser, timestamp int64) UserEvent {
	return baseUserEvent{
		insertID:  uuid.NewString(),
		timestamp: timestamp,
		user:      user,
	}
}

func NewExposureEvent(
	evaluation experiment.Evaluation,
	properties map[string]interface{},
	user user.HackleUser,
	timestamp int64,
) ExposureEvent {
	return ExposureEvent{
		UserEvent:      NewUserEvent(user, timestamp),
		Experiment:     evaluation.Experiment,
		VariationID:    evaluation.VariationID,
		VariationKey:   evaluation.VariationKey,
		DecisionReason: evaluation.Reason(),
		Properties:     properties,
	}
}

func NewTrackEvent(
	eventType model.EventType,
	event HackleEvent,
	user user.HackleUser,
	timestamp int64,
) TrackEvent {
	return TrackEvent{
		UserEvent: NewUserEvent(user, timestamp),
		EventType: eventType,
		Event:     event,
	}
}

type baseUserEvent struct {
	insertID  string
	timestamp int64
	user      user.HackleUser
}

func (e baseUserEvent) InsertID() string {
	return e.insertID
}

func (e baseUserEvent) Timestamp() int64 {
	return e.timestamp
}

func (e baseUserEvent) User() user.HackleUser {
	return e.user
}

func NewRemoteConfigEvent(
	evaluation remoteconfig.Evaluation,
	properties map[string]interface{},
	user user.HackleUser,
	timestamp int64,
) RemoteConfigEvent {
	return RemoteConfigEvent{
		UserEvent:      NewUserEvent(user, timestamp),
		Parameter:      evaluation.Parameter,
		ValueID:        evaluation.ValueID,
		DecisionReason: evaluation.Reason(),
		Properties:     properties,
	}
}

type ExposureEvent struct {
	UserEvent
	Experiment     model.Experiment
	VariationID    *int64
	VariationKey   string
	DecisionReason string
	Properties     map[string]interface{}
}

type TrackEvent struct {
	UserEvent
	EventType model.EventType
	Event     HackleEvent
}

type RemoteConfigEvent struct {
	UserEvent
	Parameter      model.RemoteConfigParameter
	ValueID        *int64
	DecisionReason string
	Properties     map[string]interface{}
}
