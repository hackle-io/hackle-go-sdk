package model

import "github.com/hackle-io/hackle-go-sdk/hackle/internal/types"

type Target struct {
	Conditions []TargetCondition
}

type TargetCondition struct {
	Key   TargetKey
	Match TargetMatch
}

type TargetKey struct {
	Type TargetKeyType
	Name string
}

type TargetKeyType string

type TargetMatch struct {
	Type      TargetMatchType
	Operator  TargetOperator
	ValueType types.ValueType
	Values    []interface{}
}

type TargetMatchType string

func (t TargetMatchType) Matches(isMatched bool) bool {
	switch t {
	case MatchTypeMatch:
		return isMatched
	case MatchTypeNotMatch:
		return !isMatched
	}
	return false
}

type TargetOperator string

const (
	TargetKeyTypeUserId         TargetKeyType = "USER_ID"
	TargetKeyTypeUserProperty   TargetKeyType = "USER_PROPERTY"
	TargetKeyTypeHackleProperty TargetKeyType = "HACKLE_PROPERTY"
	TargetKeyTypeSegment        TargetKeyType = "SEGMENT"
	TargetKeyTypeAbTest         TargetKeyType = "AB_TEST"
	TargetKeyTypeFeatureFlag    TargetKeyType = "FEATURE_FLAG"
	TargetKeyTypeEventProperty  TargetKeyType = "EVENT_PROPERTY"
)

var targetKeyTypes = map[string]TargetKeyType{
	string(TargetKeyTypeUserId):         TargetKeyTypeUserId,
	string(TargetKeyTypeUserProperty):   TargetKeyTypeUserProperty,
	string(TargetKeyTypeHackleProperty): TargetKeyTypeHackleProperty,
	string(TargetKeyTypeSegment):        TargetKeyTypeSegment,
	string(TargetKeyTypeAbTest):         TargetKeyTypeAbTest,
	string(TargetKeyTypeFeatureFlag):    TargetKeyTypeFeatureFlag,
	string(TargetKeyTypeEventProperty):  TargetKeyTypeEventProperty,
}

func TargetKeyTypeFrom(value string) (TargetKeyType, bool) {
	keyType, ok := targetKeyTypes[value]
	return keyType, ok
}

const (
	MatchTypeMatch    TargetMatchType = "MATCH"
	MatchTypeNotMatch TargetMatchType = "NOT_MATCH"
)

var targetMatchTypes = map[string]TargetMatchType{
	string(MatchTypeMatch):    MatchTypeMatch,
	string(MatchTypeNotMatch): MatchTypeNotMatch,
}

func TargetMatchTypeFrom(value string) (TargetMatchType, bool) {
	matchType, ok := targetMatchTypes[value]
	return matchType, ok
}

const (
	OperatorIn         TargetOperator = "IN"
	OperatorContains   TargetOperator = "CONTAINS"
	OperatorStartsWith TargetOperator = "STARTS_WITH"
	OperatorEndsWith   TargetOperator = "ENDS_WITH"
	OperatorGT         TargetOperator = "GT"
	OperatorGTE        TargetOperator = "GTE"
	OperatorLT         TargetOperator = "LT"
	OperatorLTE        TargetOperator = "LTE"
)

var targetOperators = map[string]TargetOperator{
	string(OperatorIn):         OperatorIn,
	string(OperatorContains):   OperatorContains,
	string(OperatorStartsWith): OperatorStartsWith,
	string(OperatorEndsWith):   OperatorEndsWith,
	string(OperatorGT):         OperatorGT,
	string(OperatorGTE):        OperatorGTE,
	string(OperatorLT):         OperatorLT,
	string(OperatorLTE):        OperatorLTE,
}

func TargetOperatorFrom(value string) (TargetOperator, bool) {
	operator, ok := targetOperators[value]
	return operator, ok
}
