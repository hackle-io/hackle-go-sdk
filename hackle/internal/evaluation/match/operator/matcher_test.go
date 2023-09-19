package operator

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMatcher(t *testing.T) {

	sut := InMatcher{}

	t.Run("string", func(t *testing.T) {
		assert.True(t, sut.StringMatches("abc", "abc"))
		assert.False(t, sut.StringMatches("abc", "abc1"))
	})

	t.Run("number", func(t *testing.T) {
		assert.True(t, sut.NumberMatches(42, 42))
		assert.True(t, sut.NumberMatches(42, 42.0))
		assert.True(t, sut.NumberMatches(42.0, 42.0))
		assert.True(t, sut.NumberMatches(42.1, 42.1))
		assert.False(t, sut.NumberMatches(42.1, 42.0))
		assert.False(t, sut.NumberMatches(43, 42))
	})

	t.Run("bool", func(t *testing.T) {
		assert.True(t, sut.BoolMatches(true, true))
		assert.True(t, sut.BoolMatches(false, false))
		assert.False(t, sut.BoolMatches(true, false))
		assert.False(t, sut.BoolMatches(false, true))
	})

	t.Run("version", func(t *testing.T) {

		assert.True(t, sut.VersionMatches(model.MustNewVersion("1.0.0"), model.MustNewVersion("1.0.0")))
		assert.False(t, sut.VersionMatches(model.MustNewVersion("1.0.0"), model.MustNewVersion("2.0.0")))
	})
}

func TestContainsMatcher(t *testing.T) {

	sut := containsMatcher{}
	t.Run("string", func(t *testing.T) {
		assert.True(t, sut.StringMatches("abc", "abc"))
		assert.True(t, sut.StringMatches("abc", "a"))
		assert.True(t, sut.StringMatches("abc", "b"))
		assert.True(t, sut.StringMatches("abc", "c"))
		assert.True(t, sut.StringMatches("abc", "ab"))
		assert.False(t, sut.StringMatches("abc", "ac"))
		assert.False(t, sut.StringMatches("a", "ab"))
	})
}

func TestMatcher(t *testing.T) {

	type s struct {
		a       string
		b       string
		matches bool
	}

	type n struct {
		a       float64
		b       float64
		matches bool
	}

	type b struct {
		a       bool
		b       bool
		matches bool
	}

	type v struct {
		a       string
		b       string
		matches bool
	}

	tests := []struct {
		name     string
		matcher  Matcher
		strings  []s
		numbers  []n
		bools    []b
		versions []v
	}{
		{
			name:    "IN",
			matcher: &InMatcher{},
			strings: []s{
				{"abc", "abc", true},
				{"abc", "abc1", false},
			},
			numbers: []n{
				{42, 42, true},
				{42.0, 42, true},
				{42.1, 42.1, true},
				{42, 42.1, false},
				{42, 43, false},
			},
			bools: []b{
				{true, true, true},
				{false, false, true},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "1.0.0", true},
				{"1.0.0", "2.0.0", false},
			},
		},
		{
			name:    "CONTAINS",
			matcher: &containsMatcher{},
			strings: []s{
				{"abc", "abc", true},
				{"abc", "a", true},
				{"abc", "b", true},
				{"abc", "c", true},
				{"abc", "ab", true},
				{"abc", "ac", false},
				{"a", "ab", false},
			},
			numbers: []n{
				{1, 1, false},
				{1, 11, false},
				{11, 1, false},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "1.0.0", false},
				{"1.0.0", "2.0.0", false},
			},
		},
		{
			name:    "STARTS_WITH",
			matcher: &startsWithMatcher{},
			strings: []s{
				{"abc", "abc", true},
				{"abc", "a", true},
				{"abc", "ab", true},
				{"abc", "b", false},
			},
			numbers: []n{
				{1, 1, false},
				{1, 11, false},
				{11, 1, false},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "1.0.0", false},
				{"1.0.0", "2.0.0", false},
			},
		},
		{
			name:    "ENDS_WITH",
			matcher: &endsWithMatcher{},
			strings: []s{
				{"abc", "abc", true},
				{"abc", "a", false},
				{"abc", "ab", false},
				{"abc", "b", false},
				{"abc", "c", true},
				{"abc", "bc", true},
			},
			numbers: []n{
				{1, 1, false},
				{1, 11, false},
				{11, 1, false},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "1.0.0", false},
				{"1.0.0", "2.0.0", false},
			},
		},
		{
			name:    "GT",
			matcher: &greaterThanMatcher{},
			strings: []s{
				{"41", "42", false},
				{"42", "42", false},
				{"43", "42", true},

				{"20230114", "20230115", false},
				{"20230115", "20230115", false},
				{"20230116", "20230115", true},

				{"2023-01-14", "2023-01-15", false},
				{"2023-01-15", "2023-01-15", false},
				{"2023-01-16", "2023-01-15", true},

				{"a", "a", false},
				{"a", "A", true},
				{"A", "a", false},
				{"aa", "a", true},
				{"a", "aa", false},
			},
			numbers: []n{
				{1, 2, false},
				{2, 2, false},
				{3, 2, true},

				{0.999, 1, false},
				{1, 1, false},
				{1.001, 1, true},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "2.0.0", false},
				{"2.0.0", "2.0.0", false},
				{"3.0.0", "2.0.0", true},
			},
		},
		{
			name:    "GTE",
			matcher: &greaterThanOrEqualToMatcher{},
			strings: []s{
				{"41", "42", false},
				{"42", "42", true},
				{"43", "42", true},

				{"20230114", "20230115", false},
				{"20230115", "20230115", true},
				{"20230116", "20230115", true},

				{"2023-01-14", "2023-01-15", false},
				{"2023-01-15", "2023-01-15", true},
				{"2023-01-16", "2023-01-15", true},

				{"a", "a", true},
				{"a", "A", true},
				{"A", "a", false},
				{"aa", "a", true},
				{"a", "aa", false},
			},
			numbers: []n{
				{1, 2, false},
				{2, 2, true},
				{3, 2, true},

				{0.999, 1, false},
				{1, 1, true},
				{1.001, 1, true},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "2.0.0", false},
				{"2.0.0", "2.0.0", true},
				{"3.0.0", "2.0.0", true},
			},
		},
		{
			name:    "LT",
			matcher: &lessThanMatcher{},
			strings: []s{
				{"41", "42", true},
				{"42", "42", false},
				{"43", "42", false},

				{"20230114", "20230115", true},
				{"20230115", "20230115", false},
				{"20230116", "20230115", false},

				{"2023-01-14", "2023-01-15", true},
				{"2023-01-15", "2023-01-15", false},
				{"2023-01-16", "2023-01-15", false},

				{"a", "a", false},
				{"a", "A", false},
				{"A", "a", true},
				{"aa", "a", false},
				{"a", "aa", true},
			},
			numbers: []n{
				{1, 2, true},
				{2, 2, false},
				{3, 2, false},

				{0.999, 1, true},
				{1, 1, false},
				{1.001, 1, false},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "2.0.0", true},
				{"2.0.0", "2.0.0", false},
				{"3.0.0", "2.0.0", false},
			},
		},
		{
			name:    "LTE",
			matcher: &lessThanOrEqualToMatcher{},
			strings: []s{
				{"41", "42", true},
				{"42", "42", true},
				{"43", "42", false},

				{"20230114", "20230115", true},
				{"20230115", "20230115", true},
				{"20230116", "20230115", false},

				{"2023-01-14", "2023-01-15", true},
				{"2023-01-15", "2023-01-15", true},
				{"2023-01-16", "2023-01-15", false},

				{"a", "a", true},
				{"a", "A", false},
				{"A", "a", true},
				{"aa", "a", false},
				{"a", "aa", true},
			},
			numbers: []n{
				{1, 2, true},
				{2, 2, true},
				{3, 2, false},

				{0.999, 1, true},
				{1, 1, true},
				{1.001, 1, false},
			},
			bools: []b{
				{true, true, false},
				{false, false, false},
				{true, false, false},
				{false, true, false},
			},
			versions: []v{
				{"1.0.0", "2.0.0", true},
				{"2.0.0", "2.0.0", true},
				{"3.0.0", "2.0.0", false},
			},
		},
	}

	for _, test := range tests {
		for _, tc := range test.strings {
			t.Run(fmt.Sprintf(test.name+" string"), func(t *testing.T) {
				assert.Equal(t, tc.matches, test.matcher.StringMatches(tc.a, tc.b))
			})
		}
		for _, tc := range test.numbers {
			t.Run(fmt.Sprintf(test.name+" number"), func(t *testing.T) {
				assert.Equal(t, tc.matches, test.matcher.NumberMatches(tc.a, tc.b))
			})
		}
		for _, tc := range test.bools {
			t.Run(fmt.Sprintf(test.name+" bool"), func(t *testing.T) {
				assert.Equal(t, tc.matches, test.matcher.BoolMatches(tc.a, tc.b))
			})
		}
		for _, tc := range test.versions {
			t.Run(fmt.Sprintf(test.name+" version"), func(t *testing.T) {
				assert.Equal(t, tc.matches, test.matcher.VersionMatches(model.MustNewVersion(tc.a), model.MustNewVersion(tc.b)))
			})
		}
	}
}
