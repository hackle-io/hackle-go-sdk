package hackle

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/clock"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/core"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/event"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/http"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/logger"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/model"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/schedule"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/workspace"
	"sync"
	"time"
)

type Client interface {
	Variation(experimentKey int64, user User) string
	VariationDetail(experimentKey int64, user User) ExperimentDecision
	IsFeatureOn(featureKey int64, user User) bool
	FeatureFlagDetail(featureKey int64, user User) FeatureFlagDecision
	RemoteConfig(user User) RemoteConfig
	Track(event Event, user User)
	Close()
}

var clients = make(map[string]Client)
var mu = &sync.Mutex{}

func NewClient(sdkKey string, config *Config) Client {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := clients[sdkKey]; ok {
		return c
	}

	c := createClient(sdkKey, config)
	clients[sdkKey] = c
	return c
}

func createClient(sdkKey string, config *Config) Client {
	if config == nil {
		config = NewConfigBuilder().Build()
	}

	sdk := model.NewSdk(sdkKey)
	scheduler := schedule.NewTickerScheduler()
	httpClient := http.NewClient(sdk, clock.System, 10*time.Second)

	httpWorkspaceFetcher := workspace.NewHttpFetcher(config.sdkUrl, sdk, httpClient)
	workspaceFetcher := workspace.NewPollingFetcher(httpWorkspaceFetcher, 10*time.Second, scheduler)

	eventDispatcher := event.NewDispatcher(config.eventUrl, httpClient)
	eventProcessor := event.NewProcessor(10000, eventDispatcher, 100, scheduler, 10*time.Second)

	c := core.New(workspaceFetcher, eventProcessor)
	userResolver := user.NewResolver()

	workspaceFetcher.Start()
	eventProcessor.Start()

	return &client{
		core:         c,
		userResolver: userResolver,
	}
}

type client struct {
	core         core.Core
	userResolver user.Resolver
}

func (c *client) Variation(experimentKey int64, user User) string {
	return c.VariationDetail(experimentKey, user).Variation()
}

func (c *client) VariationDetail(experimentKey int64, user User) ExperimentDecision {
	hackleUser, ok := c.userResolver.Resolve(user)
	if !ok {
		return decision.NewExperimentDecision("A", decision.ReasonInvalidInput, config.Empty())
	}
	d, err := c.core.Experiment(experimentKey, hackleUser, "A")
	if err != nil {
		logger.Error("Unexpected error while deciding variation for experiment[%d]. Returning control variation[A]: %v", experimentKey, err)
		return decision.NewExperimentDecision("A", decision.ReasonException, config.Empty())
	}
	return d
}

func (c *client) IsFeatureOn(featureKey int64, user User) bool {
	return c.FeatureFlagDetail(featureKey, user).IsOn()
}

func (c *client) FeatureFlagDetail(featureKey int64, user User) FeatureFlagDecision {
	hackleUser, ok := c.userResolver.Resolve(user)
	if !ok {
		return decision.NewFeatureFlagDecision(false, decision.ReasonInvalidInput, config.Empty())
	}
	d, err := c.core.FeatureFlag(featureKey, hackleUser)
	if err != nil {
		logger.Error("Unexpected error while deciding feature flag[%d]. Returning control flag[false]: %v", featureKey, err)
		return decision.NewFeatureFlagDecision(false, decision.ReasonException, config.Empty())
	}
	return d
}

func (c *client) RemoteConfig(user User) RemoteConfig {
	return newRemoteConfig(user, c.userResolver, c.core)
}

func (c *client) Track(event Event, user User) {
	hackleUser, ok := c.userResolver.Resolve(user)
	if !ok {
		return
	}
	c.core.Track(event, hackleUser)
}

func (c *client) Close() {
	c.core.Close()
}
