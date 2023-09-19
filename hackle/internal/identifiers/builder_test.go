package identifiers

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIdentifiers(t *testing.T) {
	tests := []struct {
		name        string
		identifiers map[string]string
		expected    map[string]string
	}{
		{
			name: "max identifier type length is 128",
			identifiers: map[string]string{
				strings.Repeat("a", 128): "128",
				strings.Repeat("a", 129): "129",
			},
			expected: map[string]string{
				strings.Repeat("a", 128): "128",
			},
		},
		{
			name: "max identifier value length is 512",
			identifiers: map[string]string{
				"512": strings.Repeat("a", 512),
				"513": strings.Repeat("a", 513),
			},
			expected: map[string]string{
				"512": strings.Repeat("a", 512),
			},
		},
		{
			name: "identifier type cannot be empty",
			identifiers: map[string]string{
				"": "a",
			},
			expected: map[string]string{},
		},
		{
			name: "identifier value cannot be empty",
			identifiers: map[string]string{
				"a": "",
			},
			expected: map[string]string{},
		},
	}

	for _, test := range tests {
		builder := NewBuilder()
		for k, v := range test.identifiers {
			builder.Add(k, v)
		}
		actual := builder.Build()
		assert.Equal(t, test.expected, actual)
	}

	assert.Equal(t, map[string]string{
		"a": "b", "c": "d",
	}, NewBuilder().AddAll(map[string]string{"a": "b", "c": "d"}).Build())
}
