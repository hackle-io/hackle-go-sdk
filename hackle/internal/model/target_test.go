package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTargetMatchType_Matches(t *testing.T) {
	assert.Equal(t, true, MatchTypeMatch.Matches(true))
	assert.Equal(t, false, MatchTypeMatch.Matches(false))
	assert.Equal(t, true, MatchTypeNotMatch.Matches(false))
	assert.Equal(t, false, MatchTypeNotMatch.Matches(true))
	assert.Equal(t, false, TargetMatchType("42").Matches(true))
	assert.Equal(t, false, TargetMatchType("42").Matches(false))

	test := func(value string, mat TargetMatchType, ok bool) {
		a, b := TargetMatchTypeFrom(value)
		assert.Equal(t, mat, a)
		assert.Equal(t, ok, b)
	}

	test("MATCH", MatchTypeMatch, true)
	test("NOT_MATCH", MatchTypeNotMatch, true)
	test("42", "", false)

}

func TestTargetKeyType(t *testing.T) {

	test := func(value string, keyType TargetKeyType, ok bool) {
		a, b := TargetKeyTypeFrom(value)
		assert.Equal(t, keyType, a)
		assert.Equal(t, ok, b)
	}

	test("USER_ID", TargetKeyTypeUserId, true)
	test("USER_PROPERTY", TargetKeyTypeUserProperty, true)
	test("HACKLE_PROPERTY", TargetKeyTypeHackleProperty, true)
	test("SEGMENT", TargetKeyTypeSegment, true)
	test("AB_TEST", TargetKeyTypeAbTest, true)
	test("FEATURE_FLAG", TargetKeyTypeFeatureFlag, true)
	test("EVENT_PROPERTY", TargetKeyTypeEventProperty, true)
	test("42", "", false)
}

func TestTargetOperatorFrom(t *testing.T) {
	test := func(value string, operator TargetOperator, ok bool) {
		a, b := TargetOperatorFrom(value)
		assert.Equal(t, operator, a)
		assert.Equal(t, ok, b)
	}

	test("IN", OperatorIn, true)
	test("CONTAINS", OperatorContains, true)
	test("STARTS_WITH", OperatorStartsWith, true)
	test("ENDS_WITH", OperatorEndsWith, true)
	test("GT", OperatorGT, true)
	test("GTE", OperatorGTE, true)
	test("LT", OperatorLT, true)
	test("LTE", OperatorLTE, true)
	test("42", "", false)
}
