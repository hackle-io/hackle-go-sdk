package user

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConditionMatcher(t *testing.T) {
	operatorMatcher := &mockValueOperatorMatcher{}
	matcher := NewConditionMatcher(operatorMatcher)
	assert.IsType(t, &valueResolver{}, matcher.valueResolver)
	assert.Equal(t, operatorMatcher, matcher.matcher)
}

func TestConditionMatcher_Matches(t *testing.T) {
	type fields struct {
		valueResolver *mockValueResolver
		matcher       *mockValueOperatorMatcher
	}
	type args struct {
		request   evaluator.Request
		context   evaluator.Context
		condition model.TargetCondition
	}
	type expected struct {
		matches bool
		err     error
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "when error on resolve value then return error",
			fields: fields{
				valueResolver: &mockValueResolver{returns: errors.New("value resolve error")},
				matcher:       &mockValueOperatorMatcher{returns: true},
			},
			args: args{
				request:   evaluator.SimpleRequest{},
				context:   evaluator.NewContext(),
				condition: model.TargetCondition{},
			},
			expected: expected{
				matches: false,
				err:     errors.New("value resolve error"),
			},
		},
		{
			name: "when user value is nil then return false",
			fields: fields{
				valueResolver: &mockValueResolver{returns: nil},
				matcher:       &mockValueOperatorMatcher{returns: true},
			},
			args: args{
				request:   evaluator.SimpleRequest{},
				context:   evaluator.NewContext(),
				condition: model.TargetCondition{},
			},
			expected: expected{
				matches: false,
				err:     nil,
			},
		},
		{
			name: "when user value is exist then matches user value",
			fields: fields{
				valueResolver: &mockValueResolver{returns: "42"},
				matcher:       &mockValueOperatorMatcher{returns: true},
			},
			args: args{
				request:   evaluator.SimpleRequest{},
				context:   evaluator.NewContext(),
				condition: model.TargetCondition{},
			},
			expected: expected{
				matches: true,
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &ConditionMatcher{
				valueResolver: tt.fields.valueResolver,
				matcher:       tt.fields.matcher,
			}
			matches, err := m.Matches(tt.args.request, tt.args.context, tt.args.condition)
			assert.Equal(t, tt.expected.matches, matches)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

func Test_valueResolver_resolve(t *testing.T) {
	type args struct {
		user user.HackleUser
		key  model.TargetKey
	}
	type expected struct {
		value interface{}
		ok    bool
		err   error
	}
	tests := []struct {
		name     string
		args     args
		expected expected
	}{
		{
			name: "user id present",
			args: args{
				user: user.NewHackleUserBuilder().Identifier("my_id", "42").Build(),
				key:  model.TargetKey{Type: model.TargetKeyTypeUserId, Name: "my_id"},
			},
			expected: expected{
				value: "42",
				ok:    true,
				err:   nil,
			},
		},

		{
			name: "user id absent",
			args: args{
				user: user.NewHackleUserBuilder().Identifier("my_id", "42").Build(),
				key:  model.TargetKey{Type: model.TargetKeyTypeUserId, Name: "your_id"},
			},
			expected: expected{
				value: "",
				ok:    false,
				err:   nil,
			},
		},
		{
			name: "user property present",
			args: args{
				user: user.NewHackleUserBuilder().Property("age", 42).Build(),
				key:  model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "age"},
			},
			expected: expected{
				value: 42,
				ok:    true,
				err:   nil,
			},
		},
		{
			name: "user property absent",
			args: args{
				user: user.NewHackleUserBuilder().Property("age", 42).Build(),
				key:  model.TargetKey{Type: model.TargetKeyTypeUserProperty, Name: "grade"},
			},
			expected: expected{
				value: nil,
				ok:    false,
				err:   nil,
			},
		},
		{
			name: "hackle property",
			args: args{
				user: user.NewHackleUserBuilder().Property("age", 42).Build(),
				key:  model.TargetKey{Type: model.TargetKeyTypeHackleProperty, Name: "platform"},
			},
			expected: expected{
				value: nil,
				ok:    false,
				err:   nil,
			},
		},
		{
			name: "unsupported type",
			args: args{
				user: user.NewHackleUserBuilder().Build(),
				key:  model.TargetKey{Type: model.TargetKeyTypeAbTest, Name: "42"},
			},
			expected: expected{
				value: nil,
				ok:    false,
				err:   errors.New("unsupported target key type [AB_TEST]"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &valueResolver{}
			value, ok, err := r.resolve(tt.args.user, tt.args.key)
			assert.Equal(t, tt.expected.value, value)
			assert.Equal(t, tt.expected.ok, ok)
			assert.Equal(t, tt.expected.err, err)
		})
	}
}

type mockValueResolver struct {
	returns interface{}
}

func (m *mockValueResolver) resolve(user user.HackleUser, key model.TargetKey) (interface{}, bool, error) {
	switch r := m.returns.(type) {
	case error:
		return nil, false, r
	}
	if m.returns == nil {
		return nil, false, nil
	}
	return m.returns, true, nil
}

type mockValueOperatorMatcher struct {
	returns bool
}

func (m *mockValueOperatorMatcher) Matches(userValue interface{}, match model.TargetMatch) bool {
	return m.returns
}
