package hackle

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/decision"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecordExperiment(t *testing.T) {
	metrics.AddRegistry(metrics.NewCumulativeRegistry())

	clock := &mockClock{returns: []int64{100, 142}}
	sample := metrics.NewTimerSample(clock)
	recordExperiment(sample, 123, decision.NewExperimentDecision("B", decision.ReasonTrafficAllocated, config.Empty()))

	timer := metrics.NewTimer("experiment.decision", metrics.Tags{"key": "123", "variation": "B", "reason": "TRAFFIC_ALLOCATED"})
	assert.Equal(t, int64(1), timer.Count())
	assert.Equal(t, int64(42), timer.Sum())
	assert.Equal(t, int64(42), timer.Max())
}

func TestRecordFeatureFlag(t *testing.T) {
	metrics.AddRegistry(metrics.NewCumulativeRegistry())

	clock := &mockClock{returns: []int64{100, 142}}
	sample := metrics.NewTimerSample(clock)
	recordFeatureFlag(sample, 123, decision.NewFeatureFlagDecision(true, decision.ReasonDefaultRule, config.Empty()))

	timer := metrics.NewTimer("feature.flag.decision", metrics.Tags{"key": "123", "on": "true", "reason": "DEFAULT_RULE"})
	assert.Equal(t, int64(1), timer.Count())
	assert.Equal(t, int64(42), timer.Sum())
	assert.Equal(t, int64(42), timer.Max())
}

func TestRecordRemoteConfig(t *testing.T) {
	metrics.AddRegistry(metrics.NewCumulativeRegistry())

	clock := &mockClock{returns: []int64{100, 142}}
	sample := metrics.NewTimerSample(clock)
	recordRemoteConfig(sample, "rc_key", decision.NewRemoteConfigDecision("42", decision.ReasonTargetRuleMatch))

	timer := metrics.NewTimer("remote.config.decision", metrics.Tags{"key": "rc_key", "reason": "TARGET_RULE_MATCH"})
	assert.Equal(t, int64(1), timer.Count())
	assert.Equal(t, int64(42), timer.Sum())
	assert.Equal(t, int64(42), timer.Max())
}

type mockClock struct {
	returns []int64
	count   int
}

func (m *mockClock) CurrentMillis() int64 {
	return 0
}

func (m *mockClock) Tick() int64 {
	t := m.returns[m.count]
	m.count++
	return t
}
