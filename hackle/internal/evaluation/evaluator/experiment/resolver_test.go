package experiment

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/evaluator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/mocks"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/ref"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionResolver_resolve(t *testing.T) {
	type fields struct {
		bucketer *mockBucketer
	}
	type args struct {
		request Request
		action  model.Action
	}
	type expected struct {
		variation      model.Variation
		ok             bool
		err            error
		bucketingCount int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "when unsupported type then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{},
				action: model.Action{
					Type: model.ActionType("unsupported"),
				},
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       errors.New("unsupported action type [unsupported]"),
			},
		},
		{
			name: "Variation - when variation id is nil then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					Experiment: model.Experiment{ID: 42},
				},
				action: model.Action{
					Type:        model.ActionTypeVariation,
					VariationID: nil,
				},
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       errors.New("action variation [42]"),
			},
		},
		{
			name: "Variation - when cannot found variation then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
					},
				},
				action: model.Action{
					Type:        model.ActionTypeVariation,
					VariationID: ref.Int64(1000),
				},
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       errors.New("variation [1000]"),
			},
		},
		{
			name: "Variation - when variation resolved return that variation",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
					},
				},
				action: model.Action{
					Type:        model.ActionTypeVariation,
					VariationID: ref.Int64(1001),
				},
			},
			expected: expected{
				variation: model.Variation{ID: 1001, Key: "A"},
				ok:        true,
				err:       nil,
			},
		},
		{
			name: "Bucket - when bucket is is nil then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
					},
				},
				action: model.Action{
					Type:     model.ActionTypeBucket,
					BucketID: nil,
				},
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       errors.New("action bucket [42]"),
			},
		},
		{
			name: "Bucket - when cannot found bucket then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace(),
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
					},
				},
				action: model.Action{
					Type:     model.ActionTypeBucket,
					BucketID: ref.Int64(2000),
				},
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       errors.New("bucket [2000]"),
			},
		},
		{
			name: "Bucket - when identifier not found then return nil",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(model.Bucket{ID: 2000}),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						IdentifierType: "custom_id",
					},
				},
				action: model.Action{
					Type:     model.ActionTypeBucket,
					BucketID: ref.Int64(2000),
				},
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       nil,
			},
		},
		{
			name: "Bucket - when not allocated then return nil",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(model.Bucket{ID: 2000}),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						IdentifierType: "$id",
					},
				},
				action: model.Action{
					Type:     model.ActionTypeBucket,
					BucketID: ref.Int64(2000),
				},
			},
			expected: expected{
				variation:      model.Variation{},
				ok:             false,
				err:            nil,
				bucketingCount: 1,
			},
		},
		{
			name: "Bucket - when allocated then return allocated variation",
			fields: fields{
				bucketer: &mockBucketer{returns: model.Slot{VariationID: 1002}},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(model.Bucket{ID: 2000}),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID: 42,
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						IdentifierType: "$id",
					},
				},
				action: model.Action{
					Type:     model.ActionTypeBucket,
					BucketID: ref.Int64(2000),
				},
			},
			expected: expected{
				variation:      model.Variation{ID: 1002, Key: "B"},
				ok:             true,
				err:            nil,
				bucketingCount: 1,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &actionResolver{bucketer: tc.fields.bucketer}
			variation, ok, err := sut.Resolve(tc.args.request, tc.args.action)
			assert.Equal(t, tc.expected.variation, variation)
			assert.Equal(t, tc.expected.ok, ok)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.bucketingCount, tc.fields.bucketer.count)
		})
	}
}

func TestOverrideResolver_Resolve(t *testing.T) {
	type fields struct {
		targetMatcher  *mockTargetMatcher
		actionResolver *mockActionResolver
	}
	type args struct {
		request Request
		context evaluator.Context
	}
	type expected struct {
		variation  model.Variation
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
			name: "when identifier not found then return nil",
			fields: fields{
				targetMatcher:  &mockTargetMatcher{returns: []interface{}{}},
				actionResolver: &mockActionResolver{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "custom_id",
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       nil,
			},
		},
		{
			name: "when user overridden then return overridden variation",
			fields: fields{
				targetMatcher:  &mockTargetMatcher{returns: []interface{}{}},
				actionResolver: &mockActionResolver{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						UserOverrides: map[string]int64{
							"user": 1002,
						},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				variation: model.Variation{ID: 1002, Key: "B"},
				ok:        true,
				err:       nil,
			},
		},
		{
			name: "when segment overrides is empty then return nil",
			fields: fields{
				targetMatcher:  &mockTargetMatcher{returns: []interface{}{}},
				actionResolver: &mockActionResolver{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						UserOverrides: map[string]int64{},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				variation: model.Variation{},
				ok:        false,
				err:       nil,
			},
		},
		{
			name: "when all segment override not matched then return nil",
			fields: fields{
				targetMatcher:  &mockTargetMatcher{returns: []interface{}{false, false, false, false, false}},
				actionResolver: &mockActionResolver{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						UserOverrides:    map[string]int64{},
						SegmentOverrides: []model.TargetRule{{}, {}, {}, {}, {}},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				variation:  model.Variation{},
				ok:         false,
				err:        nil,
				matchCount: 5,
			},
		},
		{
			name: "when error on segment matches then return error",
			fields: fields{
				targetMatcher:  &mockTargetMatcher{returns: []interface{}{false, false, errors.New("segment match fail"), false, false}},
				actionResolver: &mockActionResolver{returns: nil},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						UserOverrides:    map[string]int64{},
						SegmentOverrides: []model.TargetRule{{}, {}, {}, {}, {}},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				variation:  model.Variation{},
				ok:         false,
				err:        errors.New("segment match fail"),
				matchCount: 3,
			},
		},
		{
			name: "when segment overridden then return overridden variation",
			fields: fields{
				targetMatcher:  &mockTargetMatcher{returns: []interface{}{false, false, false, true, false}},
				actionResolver: &mockActionResolver{returns: model.Variation{ID: 1002, Key: "B"}},
			},
			args: args{
				request: Request{
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
						Variations: []model.Variation{
							{ID: 1001, Key: "A"},
							{ID: 1002, Key: "B"},
						},
						UserOverrides:    map[string]int64{},
						SegmentOverrides: []model.TargetRule{{}, {}, {}, {}, {}},
					},
				},
				context: evaluator.NewContext(),
			},
			expected: expected{
				variation:  model.Variation{ID: 1002, Key: "B"},
				ok:         true,
				err:        nil,
				matchCount: 4,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &overrideResolver{
				targetMatcher:  tc.fields.targetMatcher,
				actionResolver: tc.fields.actionResolver,
			}
			variation, ok, err := sut.Resolve(tc.args.request, tc.args.context)
			assert.Equal(t, tc.expected.variation, variation)
			assert.Equal(t, tc.expected.ok, ok)
			assert.Equal(t, tc.expected.err, err)
			assert.Equal(t, tc.expected.matchCount, tc.fields.targetMatcher.count)
		})
	}
}

func TestContainerResolver_IsUserInContainerGroup(t *testing.T) {
	type fields struct {
		bucketer *mockBucketer
	}
	type args struct {
		request   Request
		container model.Container
	}
	type expected struct {
		ok  bool
		err error
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected expected
	}{
		{
			name: "when identifier not found then return nil",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace(),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "custom_id",
					},
				},
				container: model.Container{
					ID:       320,
					BucketID: 1001,
					Groups:   []model.ContainerGroup{},
				},
			},
			expected: expected{
				ok:  false,
				err: nil,
			},
		},
		{
			name: "when cannot found bucket then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace(),
					user:      user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				container: model.Container{
					ID:       320,
					BucketID: 1001,
					Groups:   []model.ContainerGroup{},
				},
			},
			expected: expected{
				ok:  false,
				err: errors.New("bucket [1001]"),
			},
		},
		{
			name: "when bucket not allocated then return nil",
			fields: fields{
				bucketer: &mockBucketer{returns: nil},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(
						model.Bucket{
							ID: 1001,
						},
					),
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				container: model.Container{
					ID:       320,
					BucketID: 1001,
					Groups:   []model.ContainerGroup{},
				},
			},
			expected: expected{
				ok:  false,
				err: nil,
			},
		},
		{
			name: "when container group not found then return error",
			fields: fields{
				bucketer: &mockBucketer{returns: model.Slot{VariationID: 2001}},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(
						model.Bucket{
							ID: 1001,
						},
					),
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				container: model.Container{
					ID:       320,
					BucketID: 1001,
					Groups:   []model.ContainerGroup{},
				},
			},
			expected: expected{
				ok:  false,
				err: errors.New("container group [2001]"),
			},
		},
		{
			name: "when user not in container group then return false",
			fields: fields{
				bucketer: &mockBucketer{returns: model.Slot{VariationID: 2001}},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(
						model.Bucket{
							ID: 1001,
						},
					),
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				container: model.Container{
					ID:       320,
					BucketID: 1001,
					Groups: []model.ContainerGroup{
						{ID: 2001, Experiments: []int64{1, 2, 3}},
					},
				},
			},
			expected: expected{
				ok:  false,
				err: nil,
			},
		},
		{
			name: "when user in container group then return true",
			fields: fields{
				bucketer: &mockBucketer{returns: model.Slot{VariationID: 2001}},
			},
			args: args{
				request: Request{
					workspace: mocks.CreateWorkspace().Bucket(
						model.Bucket{
							ID: 1001,
						},
					),
					user: user.NewHackleUserBuilder().Identifier(user.IdentifierTypeID, "user").Build(),
					Experiment: model.Experiment{
						ID:             42,
						IdentifierType: "$id",
					},
				},
				container: model.Container{
					ID:       320,
					BucketID: 1001,
					Groups: []model.ContainerGroup{
						{ID: 2001, Experiments: []int64{1, 2, 3, 42}},
					},
				},
			},
			expected: expected{
				ok:  true,
				err: nil,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &containerResolver{
				bucketer: tc.fields.bucketer,
			}
			ok, err := sut.IsUserInContainerGroup(tc.args.request, tc.args.container)
			assert.Equal(t, tc.expected.ok, ok)
			assert.Equal(t, tc.expected.err, err)
		})
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
