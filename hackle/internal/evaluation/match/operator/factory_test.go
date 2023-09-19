package operator

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcherFactory(t *testing.T) {

	factory := NewMatcherFactory()

	tests := []struct {
		operator model.TargetOperator
		matcher  Matcher
	}{
		{model.OperatorIn, &InMatcher{}},
		{model.OperatorContains, &containsMatcher{}},
		{model.OperatorStartsWith, &startsWithMatcher{}},
		{model.OperatorEndsWith, &endsWithMatcher{}},
		{model.OperatorGT, &greaterThanMatcher{}},
		{model.OperatorGTE, &greaterThanOrEqualToMatcher{}},
		{model.OperatorLT, &lessThanMatcher{}},
		{model.OperatorLTE, &lessThanOrEqualToMatcher{}},
	}

	for _, tc := range tests {
		matcher, ok := factory.Get(tc.operator)
		assert.True(t, ok)
		assert.IsType(t, tc.matcher, matcher)
	}

	_, ok := factory.Get("invalid")
	assert.False(t, ok)
}
