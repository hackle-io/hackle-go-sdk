package value

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/operator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcher(t *testing.T) {

	op := &operator.InMatcher{}

	type tc struct {
		userValue  interface{}
		matchValue interface{}
		matches    bool
	}

	tests := []struct {
		name    string
		matcher Matcher
		cases   []tc
	}{
		{
			name:    "String",
			matcher: &stringMatcher{},
			cases: []tc{
				{"42", "42", true},
				{"42", 42, true},
				{42, "42", true},

				{42.42, "42.42", true},
				{"42.42", 42.42, true},
				{42.42, 42.42, true},

				{true, true, false},
				{true, 1, false},
				{"1", true, false},
			},
		},
		{
			name:    "Number",
			matcher: &numberMatcher{},
			cases: []tc{
				{42, 42, true},
				{42.42, 42.42, true},
				{0, 0, true},

				{"42", "42", true},
				{"42", 42, true},
				{42, "42", true},

				{"42.42", "42.42", true},
				{"42.42", 42.42, true},
				{42.42, "42.42", true},

				{"42.0", "42.0", true},
				{"42.0", 42.0, true},
				{42.0, "42.0", true},

				{"42a", 42, false},
				{0, "false", false},
				{0, false, false},
				{true, true, false},
			},
		},
		{
			name:    "Bool",
			matcher: &boolMatcher{},
			cases: []tc{
				{true, true, true},
				{false, false, true},
				{false, true, false},
				{true, false, false},
				{1, 1, false},
				{1, true, false},
				{0, false, false},
				{0, 0, false},
				{"true", true, false},
			},
		},
		{
			name:    "Version",
			matcher: &versionMatcher{},
			cases: []tc{
				{"1", "1", true},
				{"1", "1.0", true},
				{"1.0.0", "2.0.0", false},

				{1, "1", false},
				{"1", 1, false},
				{1, 1, false},
			},
		},
	}

	for _, test := range tests {

		t.Run("value.Matcher "+test.name, func(t *testing.T) {
			for _, tt := range test.cases {
				assert.Equal(t, tt.matches, test.matcher.Matches(op, tt.userValue, tt.matchValue))
			}

		})
	}
}
