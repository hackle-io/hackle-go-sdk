package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHackleUser(t *testing.T) {

	hackleUser := NewHackleUserBuilder().
		Identifiers(map[string]string{"type-1": "value-1"}).
		Identifier("type-2", "value-2").
		Identifier(IdentifierTypeID, "id").
		Identifier(IdentifierTypeUserID, "userID").
		Identifier(IdentifierTypeDeviceID, "deviceID").
		Properties(map[string]interface{}{"key-1": "value-1"}).
		Property("key-2", "value-2").
		Build()

	assert.Equal(t, HackleUser{
		Identifiers: map[string]string{
			"type-1":    "value-1",
			"type-2":    "value-2",
			"$id":       "id",
			"$userId":   "userID",
			"$deviceId": "deviceID",
		},
		Properties: map[string]interface{}{
			"key-1": "value-1",
			"key-2": "value-2",
		},
	}, hackleUser)
}
