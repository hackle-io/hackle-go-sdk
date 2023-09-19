package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsString(t *testing.T) {

	test := func(expected string, value interface{}) {
		actual, _ := AsString(value)
		assert.Equal(t, expected, actual)
	}

	test("string", "string")
	test("42", int(42))
	test("42", int8(42))
	test("42", int16(42))
	test("42", int32(42))
	test("42", int64(42))
	test("42", uint(42))
	test("42", uint8(42))
	test("42", uint16(42))
	test("42", uint32(42))
	test("42", uint64(42))
	test("42", float32(42))
	test("42", float64(42))
	test("", true)
}

func TestAsNumber(t *testing.T) {
	test := func(expected float64, value interface{}) {
		actual, _ := AsNumber(value)
		assert.Equal(t, expected, actual)
	}

	test(42.0, "42")
	test(42.0, "42.0")
	test(42.42, "42.42")
	test(42.0, int(42))
	test(42.0, int8(42))
	test(42.0, int16(42))
	test(42.0, int32(42))
	test(42.0, int64(42))
	test(42.0, uint(42))
	test(42.0, uint8(42))
	test(42.0, uint16(42))
	test(42.0, uint32(42))
	test(42.0, uint64(42))
	test(42.0, float32(42))
	test(42.42, float64(42.42))

	test(0, true)
	test(0, false)
}

func TestAsBool(t *testing.T) {
	test := func(expected bool, value interface{}) {
		actual, _ := AsBool(value)
		assert.Equal(t, expected, actual)
	}

	test(true, true)
	test(false, false)

	_, b := AsBool(1)
	assert.False(t, b)

	_, b2 := AsBool(0)
	assert.False(t, b2)
}

func TestAsArray(t *testing.T) {
	test := func(expected []interface{}, value interface{}) {
		actual, _ := AsArray(value)
		assert.Equal(t, expected, actual)
	}

	test([]interface{}{1, 2, 3}, []interface{}{1, 2, 3})
	test(nil, 1)
}

func TestIsNumber(t *testing.T) {
	assert.True(t, IsNumber(int(42)))
	assert.True(t, IsNumber(int8(42)))
	assert.True(t, IsNumber(int16(42)))
	assert.True(t, IsNumber(int32(42)))
	assert.True(t, IsNumber(int64(42)))
	assert.True(t, IsNumber(uint(42)))
	assert.True(t, IsNumber(uint8(42)))
	assert.True(t, IsNumber(uint16(42)))
	assert.True(t, IsNumber(uint32(42)))
	assert.True(t, IsNumber(uint64(42)))
	assert.True(t, IsNumber(float32(42)))
	assert.True(t, IsNumber(float64(42)))
	assert.False(t, IsNumber("42"))
}
