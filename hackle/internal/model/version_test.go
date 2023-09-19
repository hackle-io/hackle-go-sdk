package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func v(s string) Version {
	version, ok := NewVersion(s)
	if !ok {
		panic(s)
	}
	return version
}
func invalid(t *testing.T, v interface{}) {
	_, ok := NewVersion(v)
	assert.False(t, ok)
}

func verify(t *testing.T, version string, major int64, minor int64, patch int64, prerelease []string, build []string) {
	v, ok := NewVersion(version)
	assert.True(t, ok)

	v2 := Version{
		core:       core{major, minor, patch},
		prerelease: metadata{prerelease},
		build:      metadata{build},
	}
	assert.Equal(t, v2, v)
}

func TestNewVersion_WhenNotStringType_ThenCannotParse(t *testing.T) {
	invalid(t, nil)
	invalid(t, 42)
	invalid(t, true)
}

func TestNewVersion_WhenInvalidFormat_ThenCannotParse(t *testing.T) {
	invalid(t, "01.0.0")
	invalid(t, "1.01.0")
	invalid(t, "1.1.01")
	invalid(t, "2.x")
	invalid(t, "2.3.x")
	invalid(t, "2.3.1.4")
	invalid(t, "2.3.1*beta")
	invalid(t, "2.3.1-beta*")
	invalid(t, "2.4.1-beta_4")
}

func TestNewVersion_SemanticCoreVersion(t *testing.T) {
	verify(t, "1.0.0", 1, 0, 0, []string{}, []string{})
	verify(t, "14.165.14", 14, 165, 14, []string{}, []string{})
}

func TestNewVersion_SemanticVersionWithBuild(t *testing.T) {
	verify(t, "1.0.0-beta1", 1, 0, 0, []string{"beta1"}, []string{})
	verify(t, "1.0.0-beta.1", 1, 0, 0, []string{"beta", "1"}, []string{})
	verify(t, "1.0.0-x.y.z", 1, 0, 0, []string{"x", "y", "z"}, []string{})
}

func TestNewVersion_SemanticVersionWithPrereleaseAndBuild(t *testing.T) {
	verify(t, "1.0.0-beta.1+build.2", 1, 0, 0, []string{"beta", "1"}, []string{"build", "2"})
}

func TestNewVersion_WhenMinorOrPatchVersionIsMissing_ThenFillWithZero(t *testing.T) {
	verify(t, "15", 15, 0, 0, []string{}, []string{})
	verify(t, "15.143", 15, 143, 0, []string{}, []string{})
	verify(t, "15-x.y.z", 15, 0, 0, []string{"x", "y", "z"}, []string{})
	verify(t, "15-x.y.z+a.b.c", 15, 0, 0, []string{"x", "y", "z"}, []string{"a", "b", "c"})
}

func TestVersion_Compare(t *testing.T) {

	// core
	assert.True(t, v("2.3.4").Equals(v("2.3.4")))

	// core + prerelease
	assert.True(t, v("2.3.4-beta.1").Equals(v("2.3.4-beta.1")))
	assert.False(t, v("2.3.4-beta.1").Equals(v("2.3.4-beta.2")))

	// build
	assert.True(t, v("2.3.4+build.111").Equals(v("2.3.4+build.222")))
	assert.True(t, v("2.3.4-beta+build.111").Equals(v("2.3.4-beta+build.222")))

	// major
	assert.True(t, v("4.5.7").GreaterThan(v("3.5.7")))
	assert.True(t, v("2.5.7").LessThan(v("3.5.7")))

	// minor
	assert.True(t, v("3.6.7").GreaterThan(v("3.5.7")))
	assert.True(t, v("3.4.7").LessThan(v("3.5.7")))

	// patch
	assert.True(t, v("3.5.8").GreaterThan(v("3.5.7")))
	assert.True(t, v("3.5.6").LessThan(v("3.5.7")))

	// prerelease(numeric)
	assert.True(t, v("3.5.7-1").LessThan(v("3.5.7-2")))
	assert.True(t, v("3.5.7-1.1").LessThan(v("3.5.7-1.2")))
	assert.True(t, v("3.5.7-11").GreaterThan(v("3.5.7-1")))

	// prerelease (alphabetic)
	assert.True(t, v("3.5.7-a").Equals(v("3.5.7-a")))
	assert.True(t, v("3.5.7-a").LessThan(v("3.5.7-b")))
	assert.True(t, v("3.5.7-az").GreaterThan(v("3.5.7-ab")))

	// prerelease (alphanumeric)
	assert.True(t, v("3.5.7-9").LessThan(v("3.5.7-a")))
	assert.True(t, v("3.5.7-9").LessThan(v("3.5.7-a-9")))
	assert.True(t, v("3.5.7-beta").GreaterThan(v("3.5.7-1")))
	assert.True(t, v("3.5.7-1beta").GreaterThan(v("3.5.7-1")))

	// prerelease (length)
	assert.True(t, v("3.5.7-alpha").LessThan(v("3.5.7-alpha.1")))
	assert.True(t, v("3.5.7-1.2.3").LessThan(v("3.5.7-1.2.3.4")))
}
