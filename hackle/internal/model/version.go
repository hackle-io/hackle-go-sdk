package model

import (
	"regexp"
	"strconv"
	"strings"
)

var versionRegex = regexp.MustCompile(`^(0|[1-9]\d*)(?:\.(0|[1-9]\d*))?(?:\.(0|[1-9]\d*))?(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
var emptyMetadata = metadata{[]string{}}

type Version struct {
	core       core
	prerelease metadata
	build      metadata
}

type core struct {
	major int64
	minor int64
	patch int64
}

type metadata struct {
	identifiers []string
}

func NewVersion(v interface{}) (Version, bool) {
	s, ok := v.(string)
	if !ok {
		return Version{}, false
	}
	m := versionRegex.FindStringSubmatch(s)
	if m == nil {
		return Version{}, false
	}
	var err error

	var major, minor, patch int64
	major, err = strconv.ParseInt(m[1], 10, 64)
	if err != nil {
		return Version{}, false
	}

	if m[2] != "" {
		minor, err = strconv.ParseInt(m[2], 10, 64)
		if err != nil {
			return Version{}, false
		}
	}

	if m[3] != "" {
		patch, err = strconv.ParseInt(m[3], 10, 64)
		if err != nil {
			return Version{}, false
		}
	}

	return Version{
		core:       core{major, minor, patch},
		prerelease: newMetadata(m[4]),
		build:      newMetadata(m[5]),
	}, true
}

func MustNewVersion(v interface{}) Version {
	version, ok := NewVersion(v)
	if !ok {
		panic(v)
	}
	return version
}

func newMetadata(s string) metadata {
	if len(s) == 0 {
		return emptyMetadata
	}
	identifiers := strings.Split(s, ".")
	return metadata{identifiers}
}

func compare(a, b int64) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func compareIdentifier(a, b string) int {
	if a == b {
		return 0
	}

	if a == "" {
		if b != "" {
			return -1
		}
		return 1
	}

	if b == "" {
		if a != "" {
			return 1
		}
		return -1
	}

	ai, ae := strconv.ParseInt(a, 10, 64)
	bi, be := strconv.ParseInt(b, 10, 64)

	// a string, b string
	if ae != nil && be != nil {
		if a > b {
			return 1
		}
		return -1
	}

	// a string, b number
	if ae != nil {
		return 1
	}

	// a number, b string
	if be != nil {
		return -1
	}

	// a number, b number
	return compare(ai, bi)
}

func (v Version) compare(o Version) int {
	c := v.core.compare(o.core)
	if c != 0 {
		return c
	}
	return v.prerelease.compare(o.prerelease)
}

func (c core) compare(o core) int {
	major := compare(c.major, o.major)
	if major != 0 {
		return major
	}

	minor := compare(c.minor, o.minor)
	if minor != 0 {
		return minor
	}

	return compare(c.patch, o.patch)
}

func (m metadata) compare(o metadata) int {
	aIDs := m.identifiers
	bIDs := o.identifiers

	aLen := len(aIDs)
	bLen := len(bIDs)

	l := aLen
	if bLen < aLen {
		l = bLen
	}

	for i := 0; i < l; i++ {
		c := compareIdentifier(aIDs[i], bIDs[i])
		if c != 0 {
			return c
		}
	}
	return compare(int64(aLen), int64(bLen))
}

func (v Version) Equals(o Version) bool {
	return v.compare(o) == 0
}

func (v Version) GreaterThan(o Version) bool {
	return v.compare(o) > 0
}

func (v Version) GreaterThanOrEqual(o Version) bool {
	return v.compare(o) >= 0
}

func (v Version) LessThan(o Version) bool {
	return v.compare(o) < 0
}

func (v Version) LessThanOrEqual(o Version) bool {
	return v.compare(o) <= 0
}
