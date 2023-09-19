package event

import "github.com/hackle-io/hackle-go-sdk/hackle/internal/user"

type PayloadDTO struct {
	ExposureEvents     []ExposureEventDTO     `json:"exposureEvents"`
	TrackEvents        []TrackEventDTO        `json:"trackEvents"`
	RemoteConfigEvents []RemoteConfigEventDTO `json:"remoteConfigEvents"`
}

type ExposureEventDTO struct {
	InsertID          string                 `json:"insertId"`
	Timestamp         int64                  `json:"timestamp"`
	UserID            *string                `json:"userId"`
	Identifiers       map[string]string      `json:"identifiers"`
	UserProperties    map[string]interface{} `json:"userProperties"`
	HackleProperties  map[string]interface{} `json:"hackleProperties"`
	ExperimentID      int64                  `json:"experimentId"`
	ExperimentKey     int64                  `json:"experimentKey"`
	ExperimentType    string                 `json:"experimentType"`
	ExperimentVersion int                    `json:"experimentVersion"`
	VariationID       *int64                 `json:"variationId"`
	VariationKey      string                 `json:"variationKey"`
	DecisionReason    string                 `json:"decisionReason"`
	Properties        map[string]interface{} `json:"properties"`
}

type TrackEventDTO struct {
	InsertID         string                 `json:"insertId"`
	Timestamp        int64                  `json:"timestamp"`
	UserID           *string                `json:"userId"`
	Identifiers      map[string]string      `json:"identifiers"`
	UserProperties   map[string]interface{} `json:"userProperties"`
	HackleProperties map[string]interface{} `json:"hackleProperties"`
	EventTypeID      int64                  `json:"eventTypeId"`
	EventTypeKey     string                 `json:"eventTypeKey"`
	Value            float64                `json:"value"`
	Properties       map[string]interface{} `json:"properties"`
}

type RemoteConfigEventDTO struct {
	InsertID         string                 `json:"insertId"`
	Timestamp        int64                  `json:"timestamp"`
	UserID           *string                `json:"userId"`
	Identifiers      map[string]string      `json:"identifiers"`
	UserProperties   map[string]interface{} `json:"userProperties"`
	HackleProperties map[string]interface{} `json:"hackleProperties"`
	ParameterID      int64                  `json:"parameterId"`
	ParameterKey     string                 `json:"parameterKey"`
	ParameterType    string                 `json:"parameterType"`
	DecisionReason   string                 `json:"decisionReason"`
	ValueID          *int64                 `json:"valueId"`
	Properties       map[string]interface{} `json:"properties"`
}

func NewPayloadDTO(userEvents []UserEvent) PayloadDTO {
	exposures := make([]ExposureEventDTO, 0)
	tracks := make([]TrackEventDTO, 0)
	remoteConfigs := make([]RemoteConfigEventDTO, 0)

	for _, it := range userEvents {
		switch event := it.(type) {
		case ExposureEvent:
			exposures = append(exposures, NewExposureEventDTO(event))
			break
		case TrackEvent:
			tracks = append(tracks, NewTrackEventDTO(event))
			break
		case RemoteConfigEvent:
			remoteConfigs = append(remoteConfigs, NewRemoteConfigEventDTO(event))
		}
	}
	return PayloadDTO{
		ExposureEvents:     exposures,
		TrackEvents:        tracks,
		RemoteConfigEvents: remoteConfigs,
	}
}

func NewExposureEventDTO(event ExposureEvent) ExposureEventDTO {
	u := event.User()
	e := event.Experiment
	return ExposureEventDTO{
		InsertID:          event.InsertID(),
		Timestamp:         event.Timestamp(),
		UserID:            u.GetIdentifier(user.IdentifierTypeID),
		Identifiers:       u.Identifiers,
		UserProperties:    u.Properties,
		HackleProperties:  map[string]interface{}{},
		ExperimentID:      e.ID,
		ExperimentKey:     e.Key,
		ExperimentType:    string(e.Type),
		ExperimentVersion: e.Version,
		VariationID:       event.VariationID,
		VariationKey:      event.VariationKey,
		DecisionReason:    event.DecisionReason,
		Properties:        event.Properties,
	}
}

func NewTrackEventDTO(event TrackEvent) TrackEventDTO {
	u := event.User()
	return TrackEventDTO{
		InsertID:         event.InsertID(),
		Timestamp:        event.Timestamp(),
		UserID:           u.GetIdentifier(user.IdentifierTypeID),
		Identifiers:      u.Identifiers,
		UserProperties:   u.Properties,
		HackleProperties: map[string]interface{}{},
		EventTypeID:      event.EventType.ID,
		EventTypeKey:     event.EventType.Key,
		Value:            event.Event.Value(),
		Properties:       event.Event.Properties(),
	}
}

func NewRemoteConfigEventDTO(event RemoteConfigEvent) RemoteConfigEventDTO {
	u := event.User()
	p := event.Parameter
	return RemoteConfigEventDTO{
		InsertID:         event.InsertID(),
		Timestamp:        event.Timestamp(),
		UserID:           u.GetIdentifier(user.IdentifierTypeID),
		Identifiers:      u.Identifiers,
		UserProperties:   u.Properties,
		HackleProperties: map[string]interface{}{},
		ParameterID:      p.ID,
		ParameterKey:     p.Key,
		ParameterType:    string(p.Type),
		DecisionReason:   event.DecisionReason,
		ValueID:          event.ValueID,
		Properties:       event.Properties,
	}
}
