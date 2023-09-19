package hackle

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/identifiers"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"
)

type User struct {
	id          string
	userID      string
	deviceID    string
	identifiers map[string]string
	properties  map[string]interface{}
}

func (u User) ID() string {
	return u.id
}

func (u User) UserID() string {
	return u.userID
}

func (u User) DeviceID() string {
	return u.deviceID
}

func (u User) Identifiers() map[string]string {
	return u.identifiers
}

func (u User) Properties() map[string]interface{} {
	return u.properties
}

type UserBuilder struct {
	id          string
	userID      string
	deviceID    string
	identifiers *identifiers.Builder
	properties  *properties.Builder
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		identifiers: identifiers.NewBuilder(),
		properties:  properties.NewBuilder(),
	}
}

func (b *UserBuilder) ID(id string) *UserBuilder {
	b.id = id
	return b
}

func (b *UserBuilder) UserId(userID string) *UserBuilder {
	b.userID = userID
	return b
}

func (b *UserBuilder) DeviceID(deviceID string) *UserBuilder {
	b.deviceID = deviceID
	return b
}

func (b *UserBuilder) Identifier(identifierType string, identifierValue string) *UserBuilder {
	b.identifiers.Add(identifierType, identifierValue)
	return b
}

func (b *UserBuilder) Identifiers(identifiers map[string]string) *UserBuilder {
	b.identifiers.AddAll(identifiers)
	return b
}

func (b *UserBuilder) Property(key string, value interface{}) *UserBuilder {
	b.properties.Add(key, value)
	return b
}

func (b *UserBuilder) Properties(properties map[string]interface{}) *UserBuilder {
	b.properties.AddAll(properties)
	return b
}

func (b *UserBuilder) Build() User {
	return User{
		id:          b.id,
		userID:      b.userID,
		deviceID:    b.deviceID,
		identifiers: b.identifiers.Build(),
		properties:  b.properties.Build(),
	}
}
