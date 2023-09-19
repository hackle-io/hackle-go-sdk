package decision

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewExperimentDecision(t *testing.T) {
	cfg := config.New(map[string]interface{}{"A": "B"})
	decision := NewExperimentDecision("C", ReasonTrafficAllocated, cfg)
	assert.Equal(t, "C", decision.Variation())
	assert.Equal(t, "TRAFFIC_ALLOCATED", decision.Reason())
	assert.Equal(t, cfg, decision.Config)
	assert.Contains(t, decision.String(), "ExperimentDecision")
}

func TestNewFeatureFlagDecision(t *testing.T) {
	cfg := config.New(map[string]interface{}{"A": "B"})
	decision := NewFeatureFlagDecision(true, ReasonDefaultRule, cfg)
	assert.Equal(t, true, decision.IsOn())
	assert.Equal(t, "DEFAULT_RULE", decision.Reason())
	assert.Equal(t, cfg, decision.Config)
	assert.Contains(t, decision.String(), "FeatureFlagDecision")
}

func TestNewRemoteConfigDecision(t *testing.T) {
	decision := NewRemoteConfigDecision("42", ReasonTargetRuleMatch)
	assert.Equal(t, "42", decision.Value())
	assert.Equal(t, "TARGET_RULE_MATCH", decision.Reason())
	assert.Contains(t, decision.String(), "RemoteConfigDecision")
}
