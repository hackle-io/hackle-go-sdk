package hackle

import (
	"errors"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/event"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/types"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/user"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestNewClient(t *testing.T) {

	t.Run("config", func(t *testing.T) {
		cfg := NewConfigBuilder().
			SdkUrl("sdk_url").
			EventUrl("event_url").
			MonitoringUrl("monitoring_url").
			Build()
		assert.IsType(t, &client{}, NewClient("SDK_KEY", cfg))
	})

	t.Run("create once", func(t *testing.T) {
		cfg := NewConfigBuilder().
			SdkUrl("sdk_url").
			EventUrl("event_url").
			MonitoringUrl("monitoring_url").
			Build()

		wg := sync.WaitGroup{}
		clients := make([]Client, 100)
		for i := 0; i < 100; i++ {
			wg.Add(1)
			i := i
			go func() {
				clients[i] = NewClient("KEY", cfg)
				wg.Done()
			}()
		}
		wg.Wait()

		for i := 0; i < 100; i++ {
			assert.Equal(t, clients[0], clients[i])
		}
	})
}

func Test_client_Variation(t *testing.T) {
	type fields struct {
		core         *mockCore
		userResolver *mockUserResolver
	}
	type args struct {
		experimentKey int64
		user          User
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		variation string
	}{
		{
			name: "when user not resolved then return control variation",
			fields: fields{
				core:         &mockCore{experiment: nil},
				userResolver: &mockUserResolver{returns: nil},
			},
			args: args{
				experimentKey: 42,
				user:          User{},
			},
			variation: "A",
		},
		{
			name: "when error on core experiment then return control variation",
			fields: fields{
				core:         &mockCore{experiment: errors.New("core error")},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				experimentKey: 42,
				user:          User{},
			},
			variation: "A",
		},
		{
			name: "core decision",
			fields: fields{
				core:         &mockCore{experiment: decision.NewExperimentDecision("B", decision.ReasonTrafficAllocated, config.Empty())},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				experimentKey: 42,
				user:          User{},
			},
			variation: "B",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &client{
				core:         tc.fields.core,
				userResolver: tc.fields.userResolver,
			}
			actual := sut.Variation(tc.args.experimentKey, tc.args.user)
			assert.Equalf(t, tc.variation, actual, "Variation(%v, %v)", tc.args.experimentKey, tc.args.user)
		})
	}
}

func Test_client_VariationDetail(t *testing.T) {
	type fields struct {
		core         *mockCore
		userResolver *mockUserResolver
	}
	type args struct {
		experimentKey int64
		user          User
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		decision ExperimentDecision
	}{
		{
			name: "when user not resolved then return control variation",
			fields: fields{
				core:         &mockCore{experiment: nil},
				userResolver: &mockUserResolver{returns: nil},
			},
			args: args{
				experimentKey: 42,
				user:          User{},
			},
			decision: decision.NewExperimentDecision("A", decision.ReasonInvalidInput, config.Empty()),
		},
		{
			name: "when error on core experiment then return control variation",
			fields: fields{
				core:         &mockCore{experiment: errors.New("core error")},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				experimentKey: 42,
				user:          User{},
			},
			decision: decision.NewExperimentDecision("A", decision.ReasonException, config.Empty()),
		},
		{
			name: "core decision",
			fields: fields{
				core:         &mockCore{experiment: decision.NewExperimentDecision("B", decision.ReasonTrafficAllocated, config.Empty())},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				experimentKey: 42,
				user:          User{},
			},
			decision: decision.NewExperimentDecision("B", decision.ReasonTrafficAllocated, config.Empty()),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &client{
				core:         tc.fields.core,
				userResolver: tc.fields.userResolver,
			}
			actual := sut.VariationDetail(tc.args.experimentKey, tc.args.user)
			assert.Equalf(t, tc.decision, actual, "VariationDetail(%v, %v)", tc.args.experimentKey, tc.args.user)
		})
	}
}

func Test_client_IsFeatureOn(t *testing.T) {
	type fields struct {
		core         *mockCore
		userResolver *mockUserResolver
	}
	type args struct {
		featureKey int64
		user       User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		isOn   bool
	}{
		{
			name: "when user not resolved then return false",
			fields: fields{
				core:         &mockCore{featureFlag: nil},
				userResolver: &mockUserResolver{returns: nil},
			},
			args: args{
				featureKey: 42,
				user:       User{},
			},
			isOn: false,
		},
		{
			name: "when error on core experiment then return false",
			fields: fields{
				core:         &mockCore{featureFlag: errors.New("core error")},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				featureKey: 42,
				user:       User{},
			},
			isOn: false,
		},
		{
			name: "core decision",
			fields: fields{
				core:         &mockCore{featureFlag: decision.NewFeatureFlagDecision(true, decision.ReasonDefaultRule, config.Empty())},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				featureKey: 42,
				user:       User{},
			},
			isOn: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &client{
				core:         tc.fields.core,
				userResolver: tc.fields.userResolver,
			}
			actual := sut.IsFeatureOn(tc.args.featureKey, tc.args.user)
			assert.Equalf(t, tc.isOn, actual, "IsFeatureOn(%v, %v)", tc.args.featureKey, tc.args.user)
		})
	}
}

func Test_client_FeatureFlagDetail(t *testing.T) {
	type fields struct {
		core         *mockCore
		userResolver *mockUserResolver
	}
	type args struct {
		featureKey int64
		user       User
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		decision FeatureFlagDecision
	}{
		{
			name: "when user not resolved then return false",
			fields: fields{
				core:         &mockCore{featureFlag: nil},
				userResolver: &mockUserResolver{returns: nil},
			},
			args: args{
				featureKey: 42,
				user:       User{},
			},
			decision: decision.NewFeatureFlagDecision(false, decision.ReasonInvalidInput, config.Empty()),
		},
		{
			name: "when error on core experiment then return false",
			fields: fields{
				core:         &mockCore{featureFlag: errors.New("core error")},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				featureKey: 42,
				user:       User{},
			},
			decision: decision.NewFeatureFlagDecision(false, decision.ReasonException, config.Empty()),
		},
		{
			name: "core decision",
			fields: fields{
				core:         &mockCore{featureFlag: decision.NewFeatureFlagDecision(true, decision.ReasonDefaultRule, config.Empty())},
				userResolver: &mockUserResolver{returns: user.HackleUser{}},
			},
			args: args{
				featureKey: 42,
				user:       User{},
			},
			decision: decision.NewFeatureFlagDecision(true, decision.ReasonDefaultRule, config.Empty()),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sut := &client{
				core:         tc.fields.core,
				userResolver: tc.fields.userResolver,
			}
			actual := sut.FeatureFlagDetail(tc.args.featureKey, tc.args.user)
			assert.Equalf(t, tc.decision, actual, "FeatureFlagDetail(%v, %v)", tc.args.featureKey, tc.args.user)
		})
	}
}

func Test_client_RemoteConfig(t *testing.T) {
	t.Run("return remote config instance", func(t *testing.T) {
		sut := &client{&mockCore{}, &mockUserResolver{}}

		rc := sut.RemoteConfig(User{id: "42"})

		assert.IsType(t, &remoteConfig{}, rc)
		assert.Equal(t, User{id: "42"}, rc.(*remoteConfig).user)
	})
}

func Test_client_Track(t *testing.T) {
	t.Run("when user not resolved then do not track", func(t *testing.T) {
		core := &mockCore{}
		sut := &client{core, user.NewResolver()}
		sut.Track(NewEvent("test"), User{})
		assert.Equal(t, 0, core.trackCount)
	})

	t.Run("when user resolved then track event", func(t *testing.T) {
		core := &mockCore{}
		sut := &client{core, user.NewResolver()}
		sut.Track(NewEvent("test"), User{id: "42"})
		assert.Equal(t, 1, core.trackCount)
	})
}

func Test_client_Close(t *testing.T) {
	core := &mockCore{}
	sut := &client{core, user.NewResolver()}
	assert.Equal(t, false, core.closed)
	sut.Close()
	assert.Equal(t, true, core.closed)
}

type mockCore struct {
	experiment   interface{}
	featureFlag  interface{}
	remoteConfig interface{}
	trackCount   int
	closed       bool
}

func (m *mockCore) Experiment(experimentKey int64, user user.HackleUser, defaultVariation string) (decision.ExperimentDecision, error) {
	switch r := m.experiment.(type) {
	case decision.ExperimentDecision:
		return r, nil
	case error:
		return decision.ExperimentDecision{}, r
	}
	panic("implement me")
}

func (m *mockCore) FeatureFlag(featureKey int64, user user.HackleUser) (decision.FeatureFlagDecision, error) {
	switch r := m.featureFlag.(type) {
	case decision.FeatureFlagDecision:
		return r, nil
	case error:
		return decision.FeatureFlagDecision{}, r
	}
	panic("implement me")
}

func (m *mockCore) RemoteConfig(parameterKey string, user user.HackleUser, requiredType types.ValueType, defaultValue interface{}) (decision.RemoteConfigDecision, error) {
	switch r := m.remoteConfig.(type) {
	case decision.RemoteConfigDecision:
		return r, nil
	case error:
		return decision.RemoteConfigDecision{}, r
	}
	panic("implement me")
}

func (m *mockCore) Track(e event.HackleEvent, user user.HackleUser) {
	m.trackCount++
}

func (m *mockCore) Close() {
	m.closed = true
}

type mockUserResolver struct {
	returns interface{}
}

func (m *mockUserResolver) Resolve(user.User) (user.HackleUser, bool) {
	if u, ok := m.returns.(user.HackleUser); ok {
		return u, true
	}
	return user.HackleUser{}, false
}
