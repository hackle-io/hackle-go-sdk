package value

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/evaluation/match/operator"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
)

type Matcher interface {
	Matches(operatorMatcher operator.Matcher, userValue interface{}, matchValue interface{}) bool
}

type stringMatcher struct{}

func (m *stringMatcher) Matches(operatorMatcher operator.Matcher, userValue interface{}, matchValue interface{}) bool {
	sUserValue, ok1 := types.AsString(userValue)
	sMatchValue, ok2 := types.AsString(matchValue)
	if ok1 && ok2 {
		return operatorMatcher.StringMatches(sUserValue, sMatchValue)
	} else {
		return false
	}
}

type numberMatcher struct{}

func (m *numberMatcher) Matches(operatorMatcher operator.Matcher, userValue interface{}, matchValue interface{}) bool {
	nUserValue, ok1 := types.AsNumber(userValue)
	nMatchValue, ok2 := types.AsNumber(matchValue)
	if ok1 && ok2 {
		return operatorMatcher.NumberMatches(nUserValue, nMatchValue)
	} else {
		return false
	}
}

type boolMatcher struct{}

func (m *boolMatcher) Matches(operatorMatcher operator.Matcher, userValue interface{}, matchValue interface{}) bool {
	bUserValue, ok1 := types.AsBool(userValue)
	bMatchValue, ok2 := types.AsBool(matchValue)
	if ok1 && ok2 {
		return operatorMatcher.BoolMatches(bUserValue, bMatchValue)
	} else {
		return false
	}
}

type versionMatcher struct{}

func (m *versionMatcher) Matches(operatorMatcher operator.Matcher, userValue interface{}, matchValue interface{}) bool {
	vUserValue, ok1 := model.NewVersion(userValue)
	vMatchValue, ok2 := model.NewVersion(matchValue)
	if ok1 && ok2 {
		return operatorMatcher.VersionMatches(vUserValue, vMatchValue)
	} else {
		return false
	}
}
