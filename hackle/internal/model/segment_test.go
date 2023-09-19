package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSegmentTypeFrom(t *testing.T) {
	test := func(value string, segmentType SegmentType, ok bool) {
		a, b := SegmentTypeFrom(value)
		assert.Equal(t, segmentType, a)
		assert.Equal(t, ok, b)
	}

	test("USER_ID", SegmentTypeUserId, true)
	test("USER_PROPERTY", SegmentTypeUserProperty, true)
	test("INVALID", "", false)
}
