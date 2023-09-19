package user

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/identifiers"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"
)

type User interface {
	ID() string
	UserID() string
	DeviceID() string
	Identifiers() map[string]string
	Properties() map[string]interface{}
}

type HackleUser struct {
	Identifiers map[string]string
	Properties  map[string]interface{}
}

func (u HackleUser) GetIdentifier(identifierType string) *string {
	if identifier, ok := u.Identifiers[identifierType]; ok {
		return &identifier
	} else {
		return nil
	}
}

type HackleUserBuilder struct {
	identifiers *identifiers.Builder
	properties  *properties.Builder
}

func NewHackleUserBuilder() *HackleUserBuilder {
	return &HackleUserBuilder{
		identifiers: identifiers.NewBuilder(),
		properties:  properties.NewBuilder(),
	}
}

func (b *HackleUserBuilder) Identifier(identifierType string, identifierValue string) *HackleUserBuilder {
	b.identifiers.Add(identifierType, identifierValue)
	return b
}

func (b *HackleUserBuilder) Identifiers(identifiers map[string]string) *HackleUserBuilder {
	b.identifiers.AddAll(identifiers)
	return b
}

func (b *HackleUserBuilder) Property(key string, value interface{}) *HackleUserBuilder {
	b.properties.Add(key, value)
	return b
}

func (b *HackleUserBuilder) Properties(properties map[string]interface{}) *HackleUserBuilder {
	b.properties.AddAll(properties)
	return b
}

func (b *HackleUserBuilder) Build() HackleUser {
	return HackleUser{
		Identifiers: b.identifiers.Build(),
		Properties:  b.properties.Build(),
	}
}

const (
	IdentifierTypeID       = "$id"
	IdentifierTypeUserID   = "$userId"
	IdentifierTypeDeviceID = "$deviceId"
)
