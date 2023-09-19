package model

var (
	TargetingTypeIdentifier = TargetingType{
		supportedKeyTypes: []TargetKeyType{
			TargetKeyTypeSegment,
		},
	}

	TargetingTypeProperty = TargetingType{
		supportedKeyTypes: []TargetKeyType{
			TargetKeyTypeSegment,
			TargetKeyTypeUserProperty,
			TargetKeyTypeEventProperty,
			TargetKeyTypeHackleProperty,
			TargetKeyTypeAbTest,
			TargetKeyTypeFeatureFlag,
		},
	}

	TargetingTypeSegment = TargetingType{
		supportedKeyTypes: []TargetKeyType{
			TargetKeyTypeUserId,
			TargetKeyTypeUserProperty,
			TargetKeyTypeHackleProperty,
		},
	}
)

type TargetingType struct {
	supportedKeyTypes []TargetKeyType
}

func (t TargetingType) Supports(keyType TargetKeyType) bool {
	for _, supportedKeyType := range t.supportedKeyTypes {
		if supportedKeyType == keyType {
			return true
		}
	}
	return false
}
