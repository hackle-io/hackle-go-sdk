package user

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResolver_Resolve(t *testing.T) {
	t.Run("empty identifier", func(t *testing.T) {
		_, ok := NewResolver().Resolve(mockUser{})
		assert.Equal(t, false, ok)
	})

	t.Run("Resolved", func(t *testing.T) {
		u := mockUser{
			id:       "ID",
			userID:   "userID",
			deviceID: "deviceID",
			identifiers: map[string]string{
				"$id":       "!",
				"$userId":   "!",
				"$deviceId": "!",
				"custom_id": "custom_value",
			},
			properties: map[string]interface{}{
				"age":   42,
				"grade": "GOLD",
				"arr":   []int{1, 2, 3},
			},
		}
		hackleUser, ok := NewResolver().Resolve(u)
		assert.Equal(t, true, ok)
		assert.Equal(t, HackleUser{
			Identifiers: map[string]string{
				"$id":       "ID",
				"$userId":   "userID",
				"$deviceId": "deviceID",
				"custom_id": "custom_value",
			},
			Properties: map[string]interface{}{
				"age":   42,
				"grade": "GOLD",
				"arr":   []interface{}{1, 2, 3},
			},
		}, hackleUser)

		assert.Equal(t, (*string)(nil), hackleUser.GetIdentifier("!!"))
		assert.Equal(t, "ID", *hackleUser.GetIdentifier(IdentifierTypeID))
	})
}

type mockUser struct {
	id          string
	userID      string
	deviceID    string
	identifiers map[string]string
	properties  map[string]interface{}
}

func (u mockUser) ID() string {
	return u.id
}

func (u mockUser) UserID() string {
	return u.userID
}

func (u mockUser) DeviceID() string {
	return u.deviceID
}

func (u mockUser) Identifiers() map[string]string {
	return u.identifiers
}

func (u mockUser) Properties() map[string]interface{} {
	return u.properties
}
