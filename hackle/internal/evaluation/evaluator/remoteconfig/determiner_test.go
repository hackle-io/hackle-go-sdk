package remoteconfig

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTargetRuleDeterminer(t *testing.T) {
	determiner := NewTargetRuleDeterminer(&mockTargetMatcher{}, &mockBucketer{})
	assert.IsType(t, &targetRuleDeterminer{}, determiner)
}

func TestTargetRuleDeterminer_Determine(t *testing.T) {
	type fields struct {
		matcher *mockMatcher
	}
	type args struct {
		request Request
		context evaluator.Context
	}
	type expected struct {
		rule       model.RemoteConfigTargetRule
		ok         bool
		err        error
		matchCount int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "when target rule is empty then return nil",
			fields: fields{
				matcher: &mockMatcher{returns: []interface{}{}},
			},
			args: args{
				request: Request{
					Parameter: model.RemoteConfigParameter{
						ID:          42,
						Key:         "param",
						TargetRules: []model.RemoteConfigTargetRule{},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				rule:       model.RemoteConfigTargetRule{},
				ok:         false,
				err:        nil,
				matchCount: 0,
			},
		},
		{
			name: "when error on target rule matches then return error",
			fields: fields{
				matcher: &mockMatcher{returns: []interface{}{false, errors.New("match error")}},
			},
			args: args{
				request: Request{
					Parameter: model.RemoteConfigParameter{
						ID:  42,
						Key: "param",
						TargetRules: []model.RemoteConfigTargetRule{
							{Key: "10001"},
							{Key: "10002"},
						},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				rule:       model.RemoteConfigTargetRule{},
				ok:         false,
				err:        errors.New("match error"),
				matchCount: 2,
			},
		},
		{
			name: "when all target rule do not matched then return nil",
			fields: fields{
				matcher: &mockMatcher{returns: []interface{}{false, false, false, false, false}},
			},
			args: args{
				request: Request{
					Parameter: model.RemoteConfigParameter{
						ID:  42,
						Key: "param",
						TargetRules: []model.RemoteConfigTargetRule{
							{Key: "10001"},
							{Key: "10002"},
							{Key: "10003"},
							{Key: "10004"},
							{Key: "10005"},
						},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				rule:       model.RemoteConfigTargetRule{},
				ok:         false,
				err:        nil,
				matchCount: 5,
			},
		},
		{
			name: "when target rule matched first then return that target rule",
			fields: fields{
				matcher: &mockMatcher{returns: []interface{}{false, false, false, true, false}},
			},
			args: args{
				request: Request{
					Parameter: model.RemoteConfigParameter{
						ID:  42,
						Key: "param",
						TargetRules: []model.RemoteConfigTargetRule{
							{Key: "10001"},
							{Key: "10002"},
							{Key: "10003"},
							{Key: "10004"},
							{Key: "10005"},
						},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				rule:       model.RemoteConfigTargetRule{Key: "10004"},
				ok:         true,
				err:        nil,
				matchCount: 4,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &targetRuleDeterminer{
				matcher: tc.fields.matcher,
			}
			rule, ok, err := sut.Determine(tc.args.request, tc.args.context)
			assert.Equal(t, tc.expected.rule, rule)
			assert.Equal(t, tc.expected.ok, ok)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.matchCount, tc.fields.matcher.count)
		})
	}
}

func Test_matcher_Matches(t *testing.T) {
	type fields struct {
		targetMatcher *mockTargetMatcher
		bucketer      *mockBucketer
	}
	type args struct {
		request    Request
		context    evaluator.Context
		targetRule model.RemoteConfigTargetRule
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
			name: "when error on matches then return error",
			fields: fields{
				targetMatcher: &mockTargetMatcher{returns: []interface{}{errors.New("match error")}},
				bucketer:      &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				context: evaluator.NewContext(),
				targetRule: model.RemoteConfigTargetRule{
					BucketID: 320,
				},
			},
			expected: expected{
				matches: false,
				err:     errors.New("match error"),
			},
		},
		{
			name: "when not matched then return false",
			fields: fields{
				targetMatcher: &mockTargetMatcher{returns: []interface{}{false}},
				bucketer:      &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				context: evaluator.NewContext(),
				targetRule: model.RemoteConfigTargetRule{
					BucketID: 320,
				},
			},
			expected: expected{
				matches: false,
				err:     nil,
			},
		},
		{
			name: "when identifier not found then return false",
			fields: fields{
				targetMatcher: &mockTargetMatcher{returns: []interface{}{true}},
				bucketer:      &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "custom_id",
					},
				},
				context: evaluator.NewContext(),
				targetRule: model.RemoteConfigTargetRule{
					BucketID: 320,
				},
			},
			expected: expected{
				matches: false,
				err:     nil,
			},
		},
		{
			name: "when bucket not found then return error",
			fields: fields{
				targetMatcher: &mockTargetMatcher{returns: []interface{}{true}},
				bucketer:      &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace(),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				context: evaluator.NewContext(),
				targetRule: model.RemoteConfigTargetRule{
					BucketID: 320,
				},
			},
			expected: expected{
				matches: false,
				err:     errors.New("bucket [320]"),
			},
		},
		{
			name: "when user allocated then return true",
			fields: fields{
				targetMatcher: &mockTargetMatcher{returns: []interface{}{true}},
				bucketer:      &mockBucketer{returns: model.Slot{}},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(model.Bucket{ID: 320}),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Parameter: model.RemoteConfigParameter{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				context: evaluator.NewContext(),
				targetRule: model.RemoteConfigTargetRule{
					BucketID: 320,
				},
			},
			expected: expected{
				matches: true,
				err:     nil,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &matcher{
				targetMatcher: tc.fields.targetMatcher,
				bucketer:      tc.fields.bucketer,
			}
			matches, err := m.Matches(tc.args.request, tc.args.context, tc.args.targetRule)
			assert.Equal(t, tc.expected.matches, matches)
			assert.Equal(t, tc.expected.err, err)
		})
	}
}

type mockMatcher struct {
	returns []interface{}
	count   int
}

func (m *mockMatcher) Matches(request Request, context evaluator.Context, targetRule model.RemoteConfigTargetRule) (bool, error) {
	r := m.returns[m.count]
	m.count++

	switch rr := r.(type) {
	case bool:
		return rr, nil
	case error:
		return false, rr
	default:
		return false, nil
	}
}

type mockBucketer struct {
	returns interface{}
	count   int
}

func (m *mockBucketer) Bucketing(bucket model.Bucket, identifier string) (model.Slot, bool) {
	m.count++
	if slot, ok := m.returns.(model.Slot); ok {
		return slot, true
	}
	return model.Slot{}, false
}

type mockTargetMatcher struct {
	returns []interface{}
	count   int
}

func (m *mockTargetMatcher) Matches(request evaluator.Request, context evaluator.Context, target model.Target) (bool, error) {
	r := m.returns[m.count]
	m.count++

	switch rr := r.(type) {
	case bool:
		return rr, nil
	case error:
		return false, rr
	default:
		return false, nil
	}
}
