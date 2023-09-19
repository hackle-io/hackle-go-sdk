package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTargetingType(t *testing.T) {

	assert.Equal(t,
		[]TargetKeyType{
			TargetKeyTypeSegment,
		},
		TargetingTypeIdentifier.supportedKeyTypes,
	)

	assert.Equal(t, true, TargetingTypeIdentifier.Supports(TargetKeyTypeSegment))
	assert.Equal(t, false, TargetingTypeIdentifier.Supports(TargetKeyTypeUserId))
	assert.Equal(t, false, TargetingTypeIdentifier.Supports(TargetKeyTypeUserProperty))
	assert.Equal(t, false, TargetingTypeIdentifier.Supports(TargetKeyTypeHackleProperty))
	assert.Equal(t, false, TargetingTypeIdentifier.Supports(TargetKeyTypeAbTest))
	assert.Equal(t, false, TargetingTypeIdentifier.Supports(TargetKeyTypeFeatureFlag))
	assert.Equal(t, false, TargetingTypeIdentifier.Supports(TargetKeyTypeEventProperty))

	assert.Equal(t,
		[]TargetKeyType{
			TargetKeyTypeSegment,
			TargetKeyTypeUserProperty,
			TargetKeyTypeEventProperty,
			TargetKeyTypeHackleProperty,
			TargetKeyTypeAbTest,
			TargetKeyTypeFeatureFlag,
		},
		TargetingTypeProperty.supportedKeyTypes,
	)

	assert.Equal(t, true, TargetingTypeProperty.Supports(TargetKeyTypeSegment))
	assert.Equal(t, false, TargetingTypeProperty.Supports(TargetKeyTypeUserId))
	assert.Equal(t, true, TargetingTypeProperty.Supports(TargetKeyTypeUserProperty))
	assert.Equal(t, true, TargetingTypeProperty.Supports(TargetKeyTypeHackleProperty))
	assert.Equal(t, true, TargetingTypeProperty.Supports(TargetKeyTypeAbTest))
	assert.Equal(t, true, TargetingTypeProperty.Supports(TargetKeyTypeFeatureFlag))
	assert.Equal(t, true, TargetingTypeProperty.Supports(TargetKeyTypeEventProperty))

	assert.Equal(t,
		[]TargetKeyType{
			TargetKeyTypeUserId,
			TargetKeyTypeUserProperty,
			TargetKeyTypeHackleProperty,
		},
		TargetingTypeSegment.supportedKeyTypes,
	)
	assert.Equal(t, false, TargetingTypeSegment.Supports(TargetKeyTypeSegment))
	assert.Equal(t, true, TargetingTypeSegment.Supports(TargetKeyTypeUserId))
	assert.Equal(t, true, TargetingTypeSegment.Supports(TargetKeyTypeUserProperty))
	assert.Equal(t, true, TargetingTypeSegment.Supports(TargetKeyTypeHackleProperty))
	assert.Equal(t, false, TargetingTypeSegment.Supports(TargetKeyTypeAbTest))
	assert.Equal(t, false, TargetingTypeSegment.Supports(TargetKeyTypeFeatureFlag))
	assert.Equal(t, false, TargetingTypeSegment.Supports(TargetKeyTypeEventProperty))
}
