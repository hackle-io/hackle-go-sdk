package types

import (
	"reflect"
	"strconv"
)

func AsString(value interface{}) (string, bool) {
	switch v := value.(type) {
	case string:
		return v, true
	case int:
		return strconv.Itoa(v), true
	case int8:
		return strconv.Itoa(int(v)), true
	case int16:
		return strconv.Itoa(int(v)), true
	case int32:
		return strconv.Itoa(int(v)), true
	case int64:
		return strconv.FormatInt(v, 10), true
	case uint:
		return strconv.FormatUint(uint64(v), 10), true
	case uint8:
		return strconv.FormatUint(uint64(v), 10), true
	case uint16:
		return strconv.FormatUint(uint64(v), 10), true
	case uint32:
		return strconv.FormatUint(uint64(v), 10), true
	case uint64:
		return strconv.FormatUint(v, 10), true
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), true
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), true
	}
	return "", false
}

func AsNumber(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}

func AsBool(value interface{}) (bool, bool) {
	b, ok := value.(bool)
	return b, ok
}

func AsArray(value interface{}) ([]interface{}, bool) {
	v := reflect.ValueOf(value)
	kind := v.Kind()
	if kind == reflect.Array || kind == reflect.Slice {
		array := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			array[i] = v.Index(i).Interface()
		}
		return array, true
	}
	return nil, false
}

func IsNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64:
		return true
	case uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	default:
		return false
	}
}
