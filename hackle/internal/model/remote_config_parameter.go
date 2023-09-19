package model

import "github.com/hackle-io/hackle-go-sdk/hackle/internal/types"

type RemoteConfigParameter struct {
	ID             int64
	Key            string
	Type           types.ValueType
	IdentifierType string
	TargetRules    []RemoteConfigTargetRule
	DefaultValue   RemoteConfigValue
}

type RemoteConfigTargetRule struct {
	Key      string
	Name     string
	Target   Target
	BucketID int64
	Value    RemoteConfigValue
}
type RemoteConfigValue struct {
	ID       int64
	RawValue interface{}
}
