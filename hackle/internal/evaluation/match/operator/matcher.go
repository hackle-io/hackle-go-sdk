package operator

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"strings"
)

type Matcher interface {
	StringMatches(value string, matchValue string) bool
	NumberMatches(value float64, matchValue float64) bool
	BoolMatches(value bool, matchValue bool) bool
	VersionMatches(value model.Version, matchValue model.Version) bool
}

type InMatcher struct {
	Matcher
}

func (m *InMatcher) StringMatches(value string, matchValue string) bool {
	return value == matchValue
}

func (m *InMatcher) NumberMatches(value float64, matchValue float64) bool {
	return value == matchValue
}

func (m *InMatcher) BoolMatches(value bool, matchValue bool) bool {
	return value == matchValue
}

func (m *InMatcher) VersionMatches(value model.Version, matchValue model.Version) bool {
	return value.Equals(matchValue)
}

type containsMatcher struct {
	Matcher
}

func (m *containsMatcher) StringMatches(value string, matchValue string) bool {
	return strings.Contains(value, matchValue)
}

func (m *containsMatcher) NumberMatches(float64, float64) bool {
	return false
}

func (m *containsMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *containsMatcher) VersionMatches(model.Version, model.Version) bool {
	return false
}

type startsWithMatcher struct {
	Matcher
}

func (m *startsWithMatcher) StringMatches(value string, matchValue string) bool {
	return strings.HasPrefix(value, matchValue)
}

func (m *startsWithMatcher) NumberMatches(float64, float64) bool {
	return false
}

func (m *startsWithMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *startsWithMatcher) VersionMatches(model.Version, model.Version) bool {
	return false
}

type endsWithMatcher struct {
	Matcher
}

func (m *endsWithMatcher) StringMatches(value string, matchValue string) bool {
	return strings.HasSuffix(value, matchValue)
}

func (m *endsWithMatcher) NumberMatches(float64, float64) bool {
	return false
}

func (m *endsWithMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *endsWithMatcher) VersionMatches(model.Version, model.Version) bool {
	return false
}

type greaterThanMatcher struct {
	Matcher
}

func (m *greaterThanMatcher) StringMatches(value string, matchValue string) bool {
	return value > matchValue
}

func (m *greaterThanMatcher) NumberMatches(value float64, matchValue float64) bool {
	return value > matchValue
}

func (m *greaterThanMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *greaterThanMatcher) VersionMatches(value model.Version, matchValue model.Version) bool {
	return value.GreaterThan(matchValue)
}

type greaterThanOrEqualToMatcher struct {
	Matcher
}

func (m *greaterThanOrEqualToMatcher) StringMatches(value string, matchValue string) bool {
	return value >= matchValue
}

func (m *greaterThanOrEqualToMatcher) NumberMatches(value float64, matchValue float64) bool {
	return value >= matchValue
}

func (m *greaterThanOrEqualToMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *greaterThanOrEqualToMatcher) VersionMatches(value model.Version, matchValue model.Version) bool {
	return value.GreaterThanOrEqual(matchValue)
}

type lessThanMatcher struct {
	Matcher
}

func (m *lessThanMatcher) StringMatches(value string, matchValue string) bool {
	return value < matchValue
}

func (m *lessThanMatcher) NumberMatches(value float64, matchValue float64) bool {
	return value < matchValue
}

func (m *lessThanMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *lessThanMatcher) VersionMatches(value model.Version, matchValue model.Version) bool {
	return value.LessThan(matchValue)
}

type lessThanOrEqualToMatcher struct {
	Matcher
}

func (m *lessThanOrEqualToMatcher) StringMatches(value string, matchValue string) bool {
	return value <= matchValue
}

func (m *lessThanOrEqualToMatcher) NumberMatches(value float64, matchValue float64) bool {
	return value <= matchValue
}

func (m *lessThanOrEqualToMatcher) BoolMatches(bool, bool) bool {
	return false
}

func (m *lessThanOrEqualToMatcher) VersionMatches(value model.Version, matchValue model.Version) bool {
	return value.LessThanOrEqual(matchValue)
}
