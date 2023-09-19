package workspace

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorkspace(t *testing.T) {
	w, _ := NewFileFetcher("../../../testdata/workspace_config.json").Fetch()

	_, ok := w.GetExperiment(4)
	assert.False(t, ok)

	e5, ok := w.GetExperiment(5)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4318,
		Key:              5,
		Name:             nil,
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusDraft,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13378, "A", false, ref.Int64(1)},
			{13379, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6094)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, e5)

	e6, ok := w.GetExperiment(6)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4319,
		Key:              6,
		Name:             ref.String("experiment_6"),
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusDraft,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13380, "A", false, nil},
			{13381, "B", false, nil},
		},
		UserOverrides: map[string]int64{
			"user_1": 13380,
			"user_2": 13381,
		},
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6097)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, e6)

	e7, ok := w.GetExperiment(7)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4320,
		Key:              7,
		Name:             nil,
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          2,
		ExecutionVersion: 3,
		Variations: []model.Variation{
			{13382, "A", false, nil},
			{13383, "B", false, nil},
			{13384, "C", false, nil},
		},
		UserOverrides:    make(map[string]int64),
		SegmentOverrides: make([]model.TargetRule, 0),
		TargetAudiences: []model.Target{
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "age"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorGTE, ValueType: types.Number, Values: []interface{}{20.0}},
					},
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "age"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorLT, ValueType: types.Number, Values: []interface{}{30.0}},
					},
				},
			},
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "platform"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"android", "ios"}},
					},
				},
			},
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "membership"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.Bool, Values: []interface{}{true}},
					},
				},
			},
		},
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6100)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, e7)

	e8, ok := w.GetExperiment(8)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4321,
		Key:              8,
		Name:             nil,
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13385, "A", false, nil},
			{13386, "B", false, nil},
		},
		UserOverrides:    make(map[string]int64),
		SegmentOverrides: make([]model.TargetRule, 0),
		TargetAudiences: []model.Target{
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "address"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorContains, ValueType: types.String, Values: []interface{}{"seoul"}},
					},
				},
			},
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "name"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorStartsWith, ValueType: types.String, Values: []interface{}{"kim"}},
					},
				},
			},
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "message"},
						Match: model.TargetMatch{Type: model.MatchTypeNotMatch, Operator: model.OperatorEndsWith, ValueType: types.String, Values: []interface{}{"!"}},
					},
				},
			},
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "point"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorGT, ValueType: types.Number, Values: []interface{}{100.0}},
					},
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "point"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorLTE, ValueType: types.Number, Values: []interface{}{200.0}},
					},
				},
			},
		},
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6103)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, e8)

	e10, ok := w.GetExperiment(10)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4323,
		Key:              10,
		Name:             nil,
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusPaused,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13390, "A", false, nil},
			{13391, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6109)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, e10)

	e11, ok := w.GetExperiment(11)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4324,
		Key:              11,
		Name:             nil,
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusCompleted,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13392, "A", false, nil},
			{13393, "B", false, nil},
			{13394, "C", false, nil},
			{13395, "D", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6112)},
		ContainerID:       nil,
		WinnerVariationID: ref.Int64(13395),
	}, e11)

	f1, ok := w.GetFeatureFlag(1)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4325,
		Key:              1,
		Name:             nil,
		Type:             model.ExperimentTypeFeatureFlag,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusPaused,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13396, "A", false, nil},
			{13397, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6115)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, f1)

	f2, ok := w.GetFeatureFlag(2)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4326,
		Key:              2,
		Name:             nil,
		Type:             model.ExperimentTypeFeatureFlag,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13398, "A", false, nil},
			{13399, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6118)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, f2)

	f3, ok := w.GetFeatureFlag(3)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4327,
		Key:              3,
		Name:             nil,
		Type:             model.ExperimentTypeFeatureFlag,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13400, "A", false, nil},
			{13401, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6121)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, f3)

	f4, ok := w.GetFeatureFlag(4)
	assert.True(t, ok)
	assert.Equal(t, model.Experiment{
		ID:               4328,
		Key:              4,
		Name:             nil,
		Type:             model.ExperimentTypeFeatureFlag,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{13402, "A", false, nil},
			{13403, "B", false, nil},
		},
		UserOverrides: map[string]int64{
			"user1": 13402,
			"user2": 13403,
		},
		SegmentOverrides: make([]model.TargetRule, 0),
		TargetAudiences:  make([]model.Target, 0),
		TargetRules: []model.TargetRule{
			{
				Target: model.Target{
					Conditions: []model.TargetCondition{
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "device"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"android"}},
						},
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "version"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"1.0.0", "1.1.0"}},
						},
					},
				},
				Action: model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6125)},
			},
			{
				Target: model.Target{
					Conditions: []model.TargetCondition{
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "device"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"ios"}},
						},
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "version"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"2.0.0", "2.1.0"}},
						},
					},
				},
				Action: model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6126)},
			},
			{
				Target: model.Target{
					Conditions: []model.TargetCondition{
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "grade"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"GOLD", "SILVER"}},
						},
					},
				},
				Action: model.Action{Type: model.ActionTypeVariation, VariationID: ref.Int64(13403)},
			},
			{
				Target: model.Target{
					Conditions: []model.TargetCondition{
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "grade"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"BRONZE"}},
						},
					},
				},
				Action: model.Action{Type: model.ActionTypeVariation, VariationID: ref.Int64(13402)},
			},
		},
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6124)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, f4)

	b5823, ok := w.GetBucket(5823)
	assert.True(t, ok)
	assert.Equal(t, model.Bucket{
		ID:       5823,
		Seed:     875758774,
		SlotSize: 10000,
		Slots:    make([]model.Slot, 0),
	}, b5823)

	b5826, ok := w.GetBucket(5826)
	assert.True(t, ok)
	assert.Equal(t, model.Bucket{
		ID:       5826,
		Seed:     1616382391,
		SlotSize: 10000,
		Slots:    make([]model.Slot, 0),
	}, b5826)

	b5829, ok := w.GetBucket(5829)
	assert.True(t, ok)
	assert.Equal(t, model.Bucket{
		ID:       5829,
		Seed:     1634243589,
		SlotSize: 10000,
		Slots: []model.Slot{
			{
				StartInclusive: 0,
				EndExclusive:   667,
				VariationID:    12919,
			},
			{
				StartInclusive: 667,
				EndExclusive:   1333,
				VariationID:    12920,
			},
			{
				StartInclusive: 1333,
				EndExclusive:   2000,
				VariationID:    12921,
			},
		},
	}, b5829)

	ea, ok := w.GetEventType("a")
	assert.True(t, ok)
	assert.Equal(t, model.EventType{
		ID:  3072,
		Key: "a",
	}, ea)

	eb, ok := w.GetEventType("b")
	assert.True(t, ok)
	assert.Equal(t, model.EventType{
		ID:  3073,
		Key: "b",
	}, eb)

	s1, ok := w.GetSegment("Internal_QA")
	assert.True(t, ok)
	assert.Equal(t, model.Segment{
		ID:   34,
		Key:  "Internal_QA",
		Type: model.SegmentTypeUserId,
		Targets: []model.Target{
			{
				Conditions: []model.TargetCondition{
					{
						Key: model.TargetKey{Type: model.TargetKeyTypeUserId, Name: "$id"},
						Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{
							"68d3abc3-d19d-4838-9024-65b4e3fb2278",
							"519eddc0-ef78-4190-9221-8d8d5dffd2cc",
							"e67a62a5-af75-44e3-ae3f-4cd5b816c693",
							"54bcd6cc-9798-45d8-be6a-37f0db9da97f",
							"8346eb43-56e9-4b53-b4ec-0dbd11f44190",
							"6992e91e-0ea2-4de6-bc3e-f54abe68e35d",
							"868f6173-5850-430f-89a1-4c71c0c47696",
							"00864add-af72-4120-9a67-627d92ed4d01",
							"4e5aa440-fc26-4d03-a95d-bbdd13325c12",
							"b4f1f1d8-8a1a-4cd8-a2e4-f1c67506e19e",
							"4cd82896-0fa9-4948-a09a-c3b2e691fa03",
							"4e1d3c43-cb5f-44ce-9566-a4aa76c8f503",
							"399ad83b-3f35-466c-82db-1aa17745545b",
							"f5f39d5a-cdef-4854-80eb-bb94d3708bfa",
							"7ac770ba-c000-4f5a-a0a5-fc64f5171e62",
							"3d0781c0-3547-4b3c-970a-83fee31b80ae",
							"9d5316f3-69cf-4b7b-9223-2e6e7af64d48",
							"9fcb795e-3cfd-4e17-84d2-9da89a191baa",
							"9d2d7d8c-a0cf-458a-a7a3-1e77729bd0a3",
							"cc4bab79-bd92-4ee8-9521-91f6ddfc5964",
							"be3b69d0-2f9b-43c1-a057-20685555436b",
							"d9872fbe-bb5c-4e2a-9bf3-3a79342a7bfe",
							"b2fc7909-2753-4869-a637-5f82d4b1d692",
							"583754b2-0bf9-40de-9325-02308f678dc5",
							"7d824051-2041-4342-ab37-b619b3559bf5",
							"04fe7b7d-8220-4b04-a6fa-1f558ef081da",
							"246a5ce0-2bfe-458d-adf3-6f1c1ab8135b",
							"7cb75039-2f7a-4bc8-88a0-68185a350422",
							"84de60fb-3072-4a4c-9496-611f1aeeca6d",
							"17ddc8e6-f9b0-4039-a70c-8f83b7598685",
							"1ed2a103-7a2a-4aa5-8691-9b2a00337828",
							"ab396179-2d46-4348-a642-416e7ebd77a0"},
						},
					},
				},
			},
		},
	}, s1)

	s2, ok := w.GetSegment("test")
	assert.True(t, ok)
	assert.Equal(t, model.Segment{
		ID:      37,
		Key:     "test",
		Type:    model.SegmentTypeUserId,
		Targets: make([]model.Target, 0),
	}, s2)

	s3, ok := w.GetSegment("not_hackle")
	assert.True(t, ok)
	assert.Equal(t, model.Segment{
		ID:   81,
		Key:  "not_hackle",
		Type: model.SegmentTypeUserProperty,
		Targets: []model.Target{
			{
				Conditions: []model.TargetCondition{
					{
						Key:   model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "workspaceId"},
						Match: model.TargetMatch{Type: model.MatchTypeNotMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"22"}},
					},
				},
			},
		},
	}, s3)

	c2, ok := w.GetContainer(2)
	assert.True(t, ok)
	assert.Equal(t, model.Container{
		ID:       2,
		BucketID: 86557,
		Groups: []model.ContainerGroup{
			{ID: 3, Experiments: []int64{30767, 31073}},
		},
	}, c2)

	c25, ok := w.GetContainer(25)
	assert.True(t, ok)
	assert.Equal(t, model.Container{
		ID:       25,
		BucketID: 90597,
		Groups: []model.ContainerGroup{
			{ID: 54, Experiments: []int64{}},
		},
	}, c25)

	c34, ok := w.GetContainer(34)
	assert.True(t, ok)
	assert.Equal(t, model.Container{
		ID:       34,
		BucketID: 95105,
		Groups: []model.ContainerGroup{
			{ID: 51, Experiments: []int64{33364}},
			{ID: 52, Experiments: []int64{11205}},
			{ID: 53, Experiments: []int64{9976}},
		},
	}, c34)

	_, ok = w.GetParameterConfiguration(999)
	assert.False(t, ok)

	c1, ok := w.GetParameterConfiguration(1)
	assert.True(t, ok)
	assert.Equal(t, model.ParameterConfiguration{
		ID: 1,
		Parameters: map[string]interface{}{
			"string_key_1":  "string_value_1",
			"boolean_key_1": true,
			"int_key_1":     2147483647.0,
			"long_key_1":    92147483647.0,
			"double_key_1":  320.1523,
			"json_key_1":    "{\"json_key\": \"json_value\"}",
		},
	}, c1)

	_, ok = w.GetRemoteConfigParameter("!")
	assert.False(t, ok)

	r1, ok := w.GetRemoteConfigParameter("json_key_1")
	assert.True(t, ok)
	assert.Equal(t, model.RemoteConfigParameter{
		ID:             1,
		Key:            "json_key_1",
		Type:           types.Json,
		IdentifierType: "$id",
		TargetRules: []model.RemoteConfigTargetRule{
			{
				Key:  "29d404c5-e154-4ba2-add9-3dd261b059d6",
				Name: "target1",
				Target: model.Target{
					Conditions: []model.TargetCondition{
						{
							Key:   model.TargetKey{Type: model.TargetKeyTypeHackleProperty, Name: "condition1_key"},
							Match: model.TargetMatch{Type: model.MatchTypeMatch, Operator: model.OperatorIn, ValueType: types.String, Values: []interface{}{"value1", "value2", "value3"}},
						},
					},
				},
				BucketID: 1,
				Value: model.RemoteConfigValue{
					ID:       1,
					RawValue: "{\"json_key\": \"json_value\"}",
				},
			},
		},
		DefaultValue: model.RemoteConfigValue{
			ID:       1,
			RawValue: "{\"json_key\": \"default_value\"}",
		},
	}, r1)
}

func TestWorkspace_Invalid(t *testing.T) {
	w, _ := NewFileFetcher("../../../testdata/workspace_invalid_config.json").Fetch()

	e1, _ := w.GetExperiment(1)
	assert.Equal(t, model.Experiment{
		ID:               1,
		Key:              1,
		Name:             nil,
		Type:             model.ExperimentTypeAbTest,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{1, "A", false, nil},
			{2, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(6100)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, e1)

	_, ok := w.GetExperiment(22)
	assert.False(t, ok)

	_, ok = w.GetExperiment(23)
	assert.False(t, ok)

	f, _ := w.GetFeatureFlag(1)
	assert.Equal(t, model.Experiment{
		ID:               2,
		Key:              1,
		Name:             nil,
		Type:             model.ExperimentTypeFeatureFlag,
		IdentifierType:   "$id",
		Status:           model.ExperimentStatusRunning,
		Version:          1,
		ExecutionVersion: 1,
		Variations: []model.Variation{
			{3, "A", false, nil},
			{4, "B", false, nil},
		},
		UserOverrides:     make(map[string]int64),
		SegmentOverrides:  make([]model.TargetRule, 0),
		TargetAudiences:   make([]model.Target, 0),
		TargetRules:       make([]model.TargetRule, 0),
		DefaultRule:       model.Action{Type: model.ActionTypeBucket, BucketID: ref.Int64(5)},
		ContainerID:       nil,
		WinnerVariationID: nil,
	}, f)
}
