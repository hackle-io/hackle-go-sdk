package user

type Resolver interface {
	Resolve(user User) (HackleUser, bool)
}

func NewResolver() Resolver {
	return &resolver{}
}

type resolver struct{}

func (r *resolver) Resolve(user User) (HackleUser, bool) {
	hackleUser := NewHackleUserBuilder().
		Identifiers(user.Identifiers()).
		Identifier(IdentifierTypeID, user.ID()).
		Identifier(IdentifierTypeUserID, user.UserID()).
		Identifier(IdentifierTypeDeviceID, user.DeviceID()).
		Properties(user.Properties()).
		Build()
	if len(hackleUser.Identifiers) == 0 {
		return HackleUser{}, false
	}
	return hackleUser, true
}
