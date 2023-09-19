package config

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConfig(t *testing.T) {

	config := New(map[string]interface{}{
		"string_key":       "string_value",
		"empty_string_key": "",
		"int_key":          42,
		"zero_int_key":     0,
		"negative_int_key": -1,
		"float_key":        0.42,
		"true_bool_key":    true,
		"false_bool_key":   false,
	})

	t.Run("get string", func(t *testing.T) {
		assert.Equal(t, "string_value", config.GetString("string_key", "!!"))
		assert.Equal(t, "", config.GetString("empty_string_key", "!!"))
		assert.Equal(t, "!!", config.GetString("default", "!!"))
	})

	t.Run("get number", func(t *testing.T) {
		assert.Equal(t, 42.0, config.GetNumber("int_key", 99))
		assert.Equal(t, 0.0, config.GetNumber("zero_int_key", 99))
		assert.Equal(t, -1.0, config.GetNumber("negative_int_key", 99))
		assert.Equal(t, 0.42, config.GetNumber("float_key", 99))
		assert.Equal(t, 99.0, config.GetNumber("default", 99))
	})

	t.Run("get bool", func(t *testing.T) {
		assert.Equal(t, true, config.GetBool("true_bool_key", false))
		assert.Equal(t, false, config.GetBool("false_bool_key", true))
		assert.Equal(t, true, config.GetBool("invalid", true))
	})

	t.Run("get", func(t *testing.T) {
		v, ok := config.get(types.Version, "string_key")
		assert.Nil(t, v)
		assert.False(t, ok)

		v, ok = config.get(types.Json, "string_key")
		assert.Nil(t, v)
		assert.False(t, ok)
	})
}

func TestEmpty(t *testing.T) {
	assert.Empty(t, Empty().parameters)
	assert.Empty(t, empty.parameters)
}
