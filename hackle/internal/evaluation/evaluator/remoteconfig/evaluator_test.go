package remoteconfig

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoteConfigEvaluator_Supports(t *testing.T) {
	sut := NewEvaluator(&mockTargetRuleDeterminer{})
	assert.False(t, sut.Supports(evaluator.SimpleRequest{}))
	assert.True(t, sut.Supports(Request{}))
}

func TestRemoteConfigEvaluator_Evaluate(t *testing.T) {
	t.Run("when not remote config request return error", func(t *testing.T) {
		// given
		request := evaluator.SimpleRequest{}
		context := evaluator.NewContext()

		// when
		sut := NewEvaluator(nil)
		_, err := sut.Evaluate(request, context)

		// then
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "unsupported evaluator request")
	})
	t.Run("EvaluateRemoteConfig", func(t *testing.T) {
		// given
		request := Request{}
		context := evaluator.NewContext()

		// when
		sut := NewEvaluator(nil)
		_, err := sut.Evaluate(request, context)

		// then
		assert.Nil(t, err)
	})
}

func Test_remoteConfigEvaluator_EvaluateRemoteConfig(t *testing.T) {
	type fields struct {
		determiner *mockTargetRuleDeterminer
	}
	type args struct {
		request Request
		context evaluator.Context
	}
	type expected struct {
		evaluation Evaluation
		err        error
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "when identifier not found then return default value",
			fields: fields{
				determiner: &mockTargetRuleDeterminer{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "custom_id",
					},
					requiredType: types.String,
					defaultValue: "default",
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				evaluation: Evaluation{
					reason:            "IDENTIFIER_NOT_FOUND",
					targetEvaluations: make([]evaluator.Evaluation, 0),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "custom_id",
					},
					ValueID: nil,
					Value:   "default",
					Properties: map[string]interface{}{
						"requestValueType":    "STRING",
						"requestDefaultValue": "default",
						"returnValue":         "default",
					},
				},
				err: nil,
			},
		},
		{
			name: "when error on determine then return error",
			fields: fields{
				determiner: &mockTargetRuleDeterminer{returns: errors.New("determine error")},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
					},
					requiredType: types.String,
					defaultValue: "default",
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				evaluation: Evaluation{},
				err:        errors.New("determine error"),
			},
		},
		{
			name: "when target rule determined then return determined value",
			fields: fields{
				determiner: &mockTargetRuleDeterminer{returns: model.RemoteConfigTargetRule{
					Key:  "target_rule_key",
					Name: "target_rule_name",
					Value: model.RemoteConfigValue{
						ID:       320,
						RawValue: "target_rule_value",
					},
				}},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
						Type:           types.String,
					},
					requiredType: types.String,
					defaultValue: "default",
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				evaluation: Evaluation{
					reason:            "TARGET_RULE_MATCH",
					targetEvaluations: make([]evaluator.Evaluation, 0),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
						Type:           types.String,
					},
					ValueID: ref.Int64(320),
					Value:   "target_rule_value",
					Properties: map[string]interface{}{
						"requestValueType":    "STRING",
						"requestDefaultValue": "default",
						"returnValue":         "target_rule_value",
						"targetRuleKey":       "target_rule_key",
						"targetRuleName":      "target_rule_name",
					},
				},
				err: nil,
			},
		},
		{
			name: "when target rule not determined then return parameter default value",
			fields: fields{
				determiner: &mockTargetRuleDeterminer{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
						Type:           types.String,
						DefaultValue: model.RemoteConfigValue{
							ID:       1001,
							RawValue: "parameter_default",
						},
					},
					requiredType: types.String,
					defaultValue: "default",
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				evaluation: Evaluation{
					reason:            "DEFAULT_RULE",
					targetEvaluations: make([]evaluator.Evaluation, 0),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
						Type:           types.String,
						DefaultValue: model.RemoteConfigValue{
							ID:       1001,
							RawValue: "parameter_default",
						},
					},
					ValueID: ref.Int64(1001),
					Value:   "parameter_default",
					Properties: map[string]interface{}{
						"requestValueType":    "STRING",
						"requestDefaultValue": "default",
						"returnValue":         "parameter_default",
					},
				},
				err: nil,
			},
		},
		{
			name: "when type mismatch then return default value",
			fields: fields{
				determiner: &mockTargetRuleDeterminer{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
						Type:           types.String,
						DefaultValue: model.RemoteConfigValue{
							ID:       1001,
							RawValue: 123.45,
						},
					},
					requiredType: types.String,
					defaultValue: "default",
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				evaluation: Evaluation{
					reason:            "TYPE_MISMATCH",
					targetEvaluations: make([]evaluator.Evaluation, 0),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
						Type:           types.String,
						DefaultValue: model.RemoteConfigValue{
							ID:       1001,
							RawValue: 123.45,
						},
					},
					ValueID: nil,
					Value:   "default",
					Properties: map[string]interface{}{
						"requestValueType":    "STRING",
						"requestDefaultValue": "default",
						"returnValue":         "default",
					},
				},
				err: nil,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := NewEvaluator(tc.fields.determiner)
			evaluation, err := sut.EvaluateRemoteConfig(tc.args.request, tc.args.context)
			assert.Equal(t, tc.expected.evaluation, evaluation)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

func Test_rawValue(t *testing.T) {
	type args struct {
		valueType types.ValueType
		value     interface{}
	}
	tests := []struct {
		name     string
		args     args
		expected interface{}
	}{
		{
			name: "string - O",
			args: args{
				valueType: types.String,
				value:     "string",
			},
			expected: "string",
		},
		{
			name: "string - X",
			args: args{
				valueType: types.String,
				value:     42,
			},
			expected: nil,
		},
		{
			name: "number - O",
			args: args{
				valueType: types.Number,
				value:     42,
			},
			expected: 42.0,
		},
		{
			name: "number - X",
			args: args{
				valueType: types.Number,
				value:     "42",
			},
			expected: nil,
		},
		{
			name: "bool - O",
			args: args{
				valueType: types.Bool,
				value:     true,
			},
			expected: true,
		},
		{
			name: "number - X",
			args: args{
				valueType: types.Bool,
				value:     0,
			},
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, _ := rawValue(Request{requiredType: tt.args.valueType}, model.RemoteConfigValue{RawValue: tt.args.value})
			assert.Equal(t, tt.expected, value)
		})
	}
}

type mockTargetRuleDeterminer struct {
	returns interface{}
}

func (m *mockTargetRuleDeterminer) Determine(request Request, context evaluator.Context) (model.RemoteConfigTargetRule, bool, error) {
	switch r := m.returns.(type) {
	case model.RemoteConfigTargetRule:
		return r, true, nil
	case error:
		return model.RemoteConfigTargetRule{}, false, r
	default:
		return model.RemoteConfigTargetRule{}, false, nil
	}
}
