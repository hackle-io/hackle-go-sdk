package hackle

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/core"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_newRemoteConfig(t *testing.T) {
	rc := newRemoteConfig(User{}, user.NewResolver(), &mockCore{})
	assert.IsType(t, &remoteConfig{}, rc)
}

func Test_remoteConfig_GetString(t *testing.T) {

	type fields struct {
		user         User
		userResolver user.Resolver
		core         core.Core
	}
	type args struct {
		key          string
		defaultValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		returns string
	}{
		{
			name: "when user not resolved then return default value",
			fields: fields{
				user:         User{},
				userResolver: user.NewResolver(),
				core:         &mockCore{},
			},
			args: args{
				key:          "rc",
				defaultValue: "default",
			},
			returns: "default",
		},
		{
			name: "when error on core then return default value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: errors.New("core error")},
			},
			args: args{
				key:          "rc",
				defaultValue: "default",
			},
			returns: "default",
		},
		{
			name: "when core returned not string value then return default value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: decision.NewRemoteConfigDecision(42, decision.ReasonDefaultRule)},
			},
			args: args{
				key:          "rc",
				defaultValue: "default",
			},
			returns: "default",
		},
		{
			name: "when core returned string value then return that string value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: decision.NewRemoteConfigDecision("target", decision.ReasonDefaultRule)},
			},
			args: args{
				key:          "rc",
				defaultValue: "default",
			},
			returns: "target",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &remoteConfig{
				user:         tt.fields.user,
				userResolver: tt.fields.userResolver,
				core:         tt.fields.core,
			}
			assert.Equalf(t, tt.returns, c.GetString(tt.args.key, tt.args.defaultValue), "GetString(%v, %v)", tt.args.key, tt.args.defaultValue)
		})
	}
}

func Test_remoteConfig_GetNumber(t *testing.T) {
	type fields struct {
		user         User
		userResolver user.Resolver
		core         core.Core
	}
	type args struct {
		key          string
		defaultValue float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		returns float64
	}{
		{
			name: "when user not resolved then return default value",
			fields: fields{
				user:         User{},
				userResolver: user.NewResolver(),
				core:         &mockCore{},
			},
			args: args{
				key:          "rc",
				defaultValue: 42.0,
			},
			returns: 42.0,
		},
		{
			name: "when error on core then return default value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: errors.New("core error")},
			},
			args: args{
				key:          "rc",
				defaultValue: 42.0,
			},
			returns: 42.0,
		},
		{
			name: "when core returned not number value then return default value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: decision.NewRemoteConfigDecision("str", decision.ReasonDefaultRule)},
			},
			args: args{
				key:          "rc",
				defaultValue: 42.0,
			},
			returns: 42.0,
		},
		{
			name: "when core returned number value then return that number value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: decision.NewRemoteConfigDecision(320.0, decision.ReasonDefaultRule)},
			},
			args: args{
				key:          "rc",
				defaultValue: 42.0,
			},
			returns: 320.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &remoteConfig{
				user:         tt.fields.user,
				userResolver: tt.fields.userResolver,
				core:         tt.fields.core,
			}
			assert.Equalf(t, tt.returns, c.GetNumber(tt.args.key, tt.args.defaultValue), "GetNumber(%v, %v)", tt.args.key, tt.args.defaultValue)
		})
	}
}

func Test_remoteConfig_GetBool(t *testing.T) {
	type fields struct {
		user         User
		userResolver user.Resolver
		core         core.Core
	}
	type args struct {
		key          string
		defaultValue bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		returns bool
	}{
		{
			name: "when user not resolved then return default value",
			fields: fields{
				user:         User{},
				userResolver: user.NewResolver(),
				core:         &mockCore{},
			},
			args: args{
				key:          "rc",
				defaultValue: true,
			},
			returns: true,
		},
		{
			name: "when error on core then return default value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: errors.New("core error")},
			},
			args: args{
				key:          "rc",
				defaultValue: true,
			},
			returns: true,
		},
		{
			name: "when core returned not bool value then return default value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: decision.NewRemoteConfigDecision("false", decision.ReasonDefaultRule)},
			},
			args: args{
				key:          "rc",
				defaultValue: true,
			},
			returns: true,
		},
		{
			name: "when core returned bool value then return that bool value",
			fields: fields{
				user:         User{id: "42"},
				userResolver: user.NewResolver(),
				core:         &mockCore{remoteConfig: decision.NewRemoteConfigDecision(false, decision.ReasonDefaultRule)},
			},
			args: args{
				key:          "rc",
				defaultValue: true,
			},
			returns: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &remoteConfig{
				user:         tt.fields.user,
				userResolver: tt.fields.userResolver,
				core:         tt.fields.core,
			}
			assert.Equalf(t, tt.returns, c.GetBool(tt.args.key, tt.args.defaultValue), "GetBool(%v, %v)", tt.args.key, tt.args.defaultValue)
		})
	}
}
