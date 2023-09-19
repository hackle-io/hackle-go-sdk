package value

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcherFactory(t *testing.T) {

	factory := NewMatcherFactory()

	tests := []struct {
		valueType types.ValueType
		matcher   Matcher
	}{
		{types.String, &stringMatcher{}},
		{types.Number, &numberMatcher{}},
		{types.Bool, &boolMatcher{}},
		{types.Version, &versionMatcher{}},
		{types.Json, &stringMatcher{}},
	}

	for _, test := range tests {
		matcher, _ := factory.Get(test.valueType)
		assert.Equal(t, test.matcher, matcher)
	}

	_, ok := factory.Get("invalid")
	assert.False(t, ok)
}
