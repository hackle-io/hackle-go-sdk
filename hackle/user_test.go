package hackle

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestUser(t *testing.T) {

	t.Run("user build", func(t *testing.T) {
		user := NewUserBuilder().
			ID("id").
			UserId("userID").
			DeviceID("deviceID").
			Identifier("ID_1", "v1").
			Identifiers(map[string]string{"ID_2": "v2"}).
			Property("int_key", 42).
			Property("long_key", int64(42)).
			Property("float_key", 42.0).
			Property("boolean_key", true).
			Property("string_key", "abc 123").
			Property("nil", nil).
			Properties(map[string]interface{}{"k1": "v1", "k2": 2}).
			Properties(nil).
			Build()
		assert.Equal(t, User{
			id:       "id",
			userID:   "userID",
			deviceID: "deviceID",
			identifiers: map[string]string{
				"ID_1": "v1",
				"ID_2": "v2",
			},
			properties: map[string]interface{}{
				"int_key":     42,
				"long_key":    int64(42),
				"float_key":   42.0,
				"boolean_key": true,
				"string_key":  "abc 123",
				"k1":          "v1",
				"k2":          2,
			},
		}, user)
	})

	t.Run("max property count is 128", func(t *testing.T) {
		builder := NewUserBuilder()
		for i := 0; i < 200; i++ {
			builder.Property(strconv.Itoa(i), i)
		}
		user := builder.Build()
		assert.Equal(t, 128, len(user.Properties()))
	})

	t.Run("max property key length is 128", func(t *testing.T) {
		key128 := strings.Repeat("a", 128)
		key129 := strings.Repeat("a", 129)

		user := NewUserBuilder().Property(key128, 128).Property(key129, 129).Build()
		_, ok := user.properties[key128]
		assert.Equal(t, true, ok)
		_, ok = user.properties[key129]
		assert.Equal(t, false, ok)
	})

	t.Run("max string property value length is 1024", func(t *testing.T) {
		v1024 := strings.Repeat("a", 1024)
		v1025 := strings.Repeat("a", 1025)

		user := NewUserBuilder().Property("1024", v1024).Property("1025", v1025).Build()

		_, ok := user.properties["1024"]
		assert.Equal(t, true, ok)
		_, ok = user.properties["1025"]
		assert.Equal(t, false, ok)
	})
}
