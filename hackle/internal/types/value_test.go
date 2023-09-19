package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValueType(t *testing.T) {
	assert.Equal(t, "STRING", String.String())
	assert.Equal(t, "NUMBER", Number.String())
	assert.Equal(t, "BOOLEAN", Bool.String())
	assert.Equal(t, "VERSION", Version.String())
	assert.Equal(t, "JSON", Json.String())

	equal := func(value string, valueType ValueType, ok bool) {
		a, b := TypeFrom(value)
		assert.Equal(t, valueType, a)
		assert.Equal(t, ok, b)
	}

	equal("STRING", String, true)
	equal("NUMBER", Number, true)
	equal("BOOLEAN", Bool, true)
	equal("VERSION", Version, true)
	equal("JSON", Json, true)
	equal("INVALID", "", false)
}
