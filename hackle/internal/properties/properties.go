package properties

import "github.com/hackle-io/hackle-go-sdk/hackle/internal/types"

type Builder struct {
	properties map[string]interface{}
}

const (
	systemPropertyKeyPrefix = '$'
	maxPropertiesCount      = 128
	maxPropertyKeyLength    = 128
	maxPropertyValueLength  = 1024
)

func NewBuilder() *Builder {
	return &Builder{
		properties: make(map[string]interface{}),
	}
}

func (b *Builder) Add(key string, value interface{}) *Builder {
	if len(b.properties) >= maxPropertiesCount {
		return b
	}
	if sanitizedValue, ok := b.sanitize(key, value); ok {
		b.properties[key] = sanitizedValue
	}
	return b
}

func (b *Builder) AddAll(properties map[string]interface{}) *Builder {
	for k, v := range properties {
		b.Add(k, v)
	}
	return b
}

func (b *Builder) Build() map[string]interface{} {
	properties := make(map[string]interface{})
	for k, v := range b.properties {
		properties[k] = v
	}
	return properties
}

func (b *Builder) sanitize(key string, value interface{}) (interface{}, bool) {
	if !b.isValidKey(key) {
		return nil, false
	}
	if value == nil {
		return nil, false
	}

	if values, ok := types.AsArray(value); ok {
		array := make([]interface{}, 0)
		for _, e := range values {
			if b.isValidElement(e) {
				array = append(array, e)
			}
		}
		return array, true
	}

	if b.isValidValue(value) {
		return value, true
	}
	if key[0] == systemPropertyKeyPrefix {
		return value, true
	}
	return nil, false
}

func (b *Builder) isValidKey(key string) bool {
	if len(key) == 0 {
		return false
	}
	if len(key) > maxPropertyKeyLength {
		return false
	}
	return true
}

func (b *Builder) isValidValue(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return len(v) <= maxPropertyValueLength
	case bool:
		return true
	default:
		return types.IsNumber(v)
	}
}

func (b *Builder) isValidElement(e interface{}) bool {
	switch v := e.(type) {
	case string:
		return len(v) <= maxPropertyValueLength
	default:
		return types.IsNumber(v)
	}
}
