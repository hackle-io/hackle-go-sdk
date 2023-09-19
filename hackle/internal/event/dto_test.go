package event

import (
	"encoding/json"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestName(t *testing.T) {

	userEvents := []UserEvent{
		ExposureEvent{
			UserEvent: baseUserEvent{
				insertID:  "experiment-1",
				timestamp: 10000,
				user: user.HackleUser{
					Identifiers: map[string]string{
						"$id":       "id",
						"$deviceId": "deviceId",
					},
					Properties: map[string]interface{}{
						"age": 42.0,
					},
				},
			},
			Experiment: model.Experiment{
				ID:               1,
				Key:              2,
				Type:             model.ExperimentTypeAbTest,
				Version:          3,
				ExecutionVersion: 4,
			},
			VariationID:    ref.Int64(4),
			VariationKey:   "A",
			DecisionReason: "EXPERIMENT_DRAFT",
			Properties: map[string]interface{}{
				"$experiment_version": 5.0,
				"$execution_version":  6.0,
			},
		},
		TrackEvent{
			UserEvent: baseUserEvent{
				insertID:  "track-1",
				timestamp: 20000,
				user: user.HackleUser{
					Identifiers: map[string]string{
						"$id":       "id",
						"$deviceId": "deviceId",
					},
					Properties: map[string]interface{}{
						"age": 42.0,
					},
				},
			},
			EventType: model.EventType{
				ID:  101,
				Key: "test",
			},
			Event: event{
				key:   "test",
				value: 42.0,
				properties: map[string]interface{}{
					"a": "b",
				},
			},
		},
		RemoteConfigEvent{
			UserEvent: baseUserEvent{
				insertID:  "rc-1",
				timestamp: 30000,
				user: user.HackleUser{
					Identifiers: map[string]string{
						"$id":       "id",
						"$deviceId": "deviceId",
					},
					Properties: map[string]interface{}{
						"age": 42.0,
					},
				},
			},
			Parameter: model.RemoteConfigParameter{
				ID:   201,
				Key:  "rc_key",
				Type: types.String,
			},
			ValueID:        ref.Int64(202),
			DecisionReason: "DEFAULT_RULE",
			Properties: map[string]interface{}{
				"requestValueType":    "STRING",
				"requestDefaultValue": "default",
				"returnValue":         "return",
			},
		},
	}

	bytes, _ := ioutil.ReadFile("../../../testdata/payload.json")
	var dto2 PayloadDTO
	json.Unmarshal(bytes, &dto2)
	assert.Equal(t, dto2, NewPayloadDTO(userEvents))
}
