package properties

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestBuilder(t *testing.T) {

	t.Run("raw value valid build", func(t *testing.T) {
		assert.Equal(t, map[string]interface{}{"key1": 1}, NewBuilder().Add("key1", 1).Build())
		assert.Equal(t, map[string]interface{}{"key1": "1"}, NewBuilder().Add("key1", "1").Build())
		assert.Equal(t, map[string]interface{}{"key1": true}, NewBuilder().Add("key1", true).Build())
		assert.Equal(t, map[string]interface{}{"key1": false}, NewBuilder().Add("key1", false).Build())
	})

	t.Run("raw invalid value", func(t *testing.T) {
		assert.Equal(t, make(map[string]interface{}), NewBuilder().Add("key1", NewBuilder()).Build())
	})

	t.Run("array value", func(t *testing.T) {
		NewBuilder().
			Add("key1", []interface{}{1, 2, 3}).
			Build()

		assert.Equal(t, map[string]interface{}{"key1": []interface{}{1, 2, 3}}, NewBuilder().Add("key1", []interface{}{1, 2, 3}).Build())
		assert.Equal(t, map[string]interface{}{"key1": []interface{}{"1", "2", "3"}}, NewBuilder().Add("key1", []interface{}{"1", "2", "3"}).Build())
		assert.Equal(t, map[string]interface{}{"key1": []interface{}{"1", 2, "3"}}, NewBuilder().Add("key1", []interface{}{"1", 2, "3"}).Build())
		assert.Equal(t, map[string]interface{}{"key1": []interface{}{1, 2, 3}}, NewBuilder().Add("key1", []interface{}{1, 2, nil, 3}).Build())
		assert.Equal(t, map[string]interface{}{"key1": []interface{}{}}, NewBuilder().Add("key1", []interface{}{true, false}).Build())
		assert.Equal(t, map[string]interface{}{"key1": []interface{}{}}, NewBuilder().Add("key1", []interface{}{strings.Repeat("a", 1025)}).Build())
	})

	t.Run("max property size is 128", func(t *testing.T) {
		builder := NewBuilder()
		for i := 0; i < 128; i++ {
			builder.Add(strconv.Itoa(i), i)
		}
		properties := builder.Build()
		assert.Len(t, properties, 128)

		builder.Add("key", 42)
		properties = builder.Build()
		assert.Len(t, properties, 128)
		assert.Nil(t, properties["key"])
	})

	t.Run("max key length is 128", func(t *testing.T) {
		builder := NewBuilder()
		builder.Add(strings.Repeat("a", 128), 128)
		assert.Len(t, builder.Build(), 1)

		builder.Add(strings.Repeat("a", 129), 129)
		assert.Len(t, builder.Build(), 1)
	})

	t.Run("empty key", func(t *testing.T) {
		properties := NewBuilder().Add("", 1).Build()
		assert.Equal(t, 0, len(properties))
	})

	t.Run("AddAll", func(t *testing.T) {
		properties := map[string]interface{}{
			"k1": "v1",
			"k2": 2,
			"k3": true,
			"k4": false,
			"k5": []int{1, 2, 3},
			"k6": []string{"1", "2", "3"},
			"k7": nil,
		}
		actual := NewBuilder().
			AddAll(properties).
			Build()
		assert.Len(t, actual, 6)
		assert.Nil(t, actual["k7"])
	})

	t.Run("add system properties", func(t *testing.T) {
		properties := NewBuilder().
			Add("$set", map[string]interface{}{"age": 42}).
			Add("set", map[string]interface{}{"age": 42}).
			Build()
		assert.Equal(t, map[string]interface{}{"$set": map[string]interface{}{"age": 42}}, properties)
	})
}
