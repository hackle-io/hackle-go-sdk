package model

type Segment struct {
	ID      int64
	Key     string
	Type    SegmentType
	Targets []Target
}

type SegmentType string

const (
	SegmentTypeUserId       SegmentType = "USER_ID"
	SegmentTypeUserProperty SegmentType = "USER_PROPERTY"
)

var segmentTypes = map[string]SegmentType{
	string(SegmentTypeUserId):       SegmentTypeUserId,
	string(SegmentTypeUserProperty): SegmentTypeUserProperty,
}

func SegmentTypeFrom(value string) (SegmentType, bool) {
	segmentType, ok := segmentTypes[value]
	return segmentType, ok
}
