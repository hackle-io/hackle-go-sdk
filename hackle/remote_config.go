package hackle

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/core"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
)

type RemoteConfig interface {
	GetString(key string, defaultValue string) string
	GetNumber(key string, defaultValue float64) float64
	GetBool(key string, defaultValue bool) bool
}

func newRemoteConfig(user User, userResolve user.Resolver, core core.Core) RemoteConfig {
	return &remoteConfig{
		user:         user,
		userResolver: userResolve,
		core:         core,
	}
}

type remoteConfig struct {
	user         User
	userResolver user.Resolver
	core         core.Core
}

func (c *remoteConfig) GetString(key string, defaultValue string) string {
	value := c.get(key, defaultValue, types.String).Value()
	if s, ok := value.(string); ok {
		return s
	} else {
		return defaultValue
	}
}

func (c *remoteConfig) GetNumber(key string, defaultValue float64) float64 {
	value := c.get(key, defaultValue, types.Number).Value()
	if n, ok := value.(float64); ok {
		return n
	} else {
		return defaultValue
	}
}

func (c *remoteConfig) GetBool(key string, defaultValue bool) bool {
	value := c.get(key, defaultValue, types.Bool).Value()
	if b, ok := value.(bool); ok {
		return b
	} else {
		return defaultValue
	}
}

func (c *remoteConfig) get(key string, defaultValue interface{}, valueType types.ValueType) RemoteConfigDecision {
	sample := metrics.NewTimerSample(clock.System)
	d := c.getInternal(key, defaultValue, valueType)
	recordRemoteConfig(sample, key, d)
	return d
}

func (c *remoteConfig) getInternal(key string, defaultValue interface{}, valueType types.ValueType) RemoteConfigDecision {
	hackleUser, ok := c.userResolver.Resolve(c.user)
	if !ok {
		return decision.NewRemoteConfigDecision(defaultValue, decision.ReasonInvalidInput)
	}
	d, err := c.core.RemoteConfig(key, hackleUser, valueType, defaultValue)
	if err != nil {
		logger.Error("Unexpected exception while deciding remote config parameter[%s]. Returning default value[%v]: %v", key, defaultValue, err)
		return decision.NewRemoteConfigDecision(defaultValue, decision.ReasonException)
	}
	return d
}
