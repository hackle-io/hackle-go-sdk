package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestActionTypeFrom(t *testing.T) {

	test := func(value string, expected ActionType) {
		v, _ := ActionTypeFrom(value)
		assert.Equal(t, expected, v)
	}

	test("VARIATION", ActionTypeVariation)
	test("BUCKET", ActionTypeBucket)
	test("INVALID", "")
}
