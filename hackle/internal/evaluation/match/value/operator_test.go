package value

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_operatorMatcher_Matches(t *testing.T) {
	type args struct {
		userValue interface{}
		match     model.TargetMatch
	}
	tests := []struct {
		name    string
		args    args
		matches bool
	}{
		{
			name: "true",
			args: args{
				userValue: 3,
				match: model.TargetMatch{
					Type:      model.MatchTypeMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: true,
		},
		{
			name: "false",
			args: args{
				userValue: 4,
				match: model.TargetMatch{
					Type:      model.MatchTypeMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: false,
		},
		{
			name: "NOT_MATCH false",
			args: args{
				userValue: 3,
				match: model.TargetMatch{
					Type:      model.MatchTypeNotMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: false,
		},
		{
			name: "NOT_MATCH true",
			args: args{
				userValue: 4,
				match: model.TargetMatch{
					Type:      model.MatchTypeNotMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: true,
		},
		{
			name: "array true",
			args: args{
				userValue: []interface{}{2},
				match: model.TargetMatch{
					Type:      model.MatchTypeMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: true,
		},
		{
			name: "array false",
			args: args{
				userValue: []interface{}{4, 5, 6},
				match: model.TargetMatch{
					Type:      model.MatchTypeMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: false,
		},
		{
			name: "empty array false",
			args: args{
				userValue: []interface{}{},
				match: model.TargetMatch{
					Type:      model.MatchTypeMatch,
					Operator:  model.OperatorIn,
					ValueType: types.Number,
					Values:    []interface{}{1, 2, 3},
				},
			},
			matches: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewOperatorMatcher()
			assert.Equalf(t, tt.matches, m.Matches(tt.args.userValue, tt.args.match), "Matches(%v, %v)", tt.args.userValue, tt.args.match)
		})
	}
}
