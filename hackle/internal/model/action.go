package model

type Action struct {
	Type        ActionType
	VariationID *int64
	BucketID    *int64
}

type ActionType string

const (
	ActionTypeVariation ActionType = "VARIATION"
	ActionTypeBucket    ActionType = "BUCKET"
)

var actionTypes = map[string]ActionType{
	string(ActionTypeVariation): ActionTypeVariation,
	string(ActionTypeBucket):    ActionTypeBucket,
}

func ActionTypeFrom(value string) (ActionType, bool) {
	actionType, ok := actionTypes[value]
	return actionType, ok
}
