package config

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
)

type Config struct {
	parameters map[string]interface{}
}

var empty = Config{make(map[string]interface{})}

func New(parameters map[string]interface{}) Config {
	return Config{
		parameters: parameters,
	}
}

func Empty() Config {
	return empty
}

func (c Config) GetString(key string, defaultValue string) string {
	value, ok := c.get(types.String, key)
	if !ok {
		return defaultValue
	}
	if v, ok := value.(string); ok {
		return v
	} else {
		return defaultValue
	}
}

func (c Config) GetNumber(key string, defaultValue float64) float64 {
	value, ok := c.get(types.Number, key)
	if !ok {
		return defaultValue
	}
	if v, ok := value.(float64); ok {
		return v
	} else {
		return defaultValue
	}
}

func (c Config) GetBool(key string, defaultValue bool) bool {
	value, ok := c.get(types.Bool, key)
	if !ok {
		return defaultValue
	}
	if v, ok := value.(bool); ok {
		return v
	} else {
		return defaultValue
	}
}

func (c Config) get(valueType types.ValueType, key string) (interface{}, bool) {
	value, ok := c.parameters[key]
	if !ok {
		return nil, false
	}
	switch valueType {
	case types.String:
		if s, ok := value.(string); ok {
			return s, true
		}
	case types.Number:
		if types.IsNumber(value) {
			return types.AsNumber(value)
		}
	case types.Bool:
		if b, ok := value.(bool); ok {
			return b, true
		}
	}
	return nil, false
}
