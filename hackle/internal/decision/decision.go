package decision

import (
	"fmt"
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/config"
)

type ExperimentDecision struct {
	config.Config
	variation string
	reason    string
}

func NewExperimentDecision(variation string, reason string, config config.Config) ExperimentDecision {
	return ExperimentDecision{
		Config:    config,
		variation: variation,
		reason:    reason,
	}
}

func (d ExperimentDecision) Variation() string {
	return d.variation
}

func (d ExperimentDecision) Reason() string {
	return d.reason
}

func (d ExperimentDecision) String() string {
	return fmt.Sprintf("ExperimentDecision(variation=%s, reason=%s, config=%s)", d.Variation(), d.Reason(), d.Config)
}

type FeatureFlagDecision struct {
	config.Config
	isOn   bool
	reason string
}

func NewFeatureFlagDecision(isOn bool, reason string, config config.Config) FeatureFlagDecision {
	return FeatureFlagDecision{
		Config: config,
		isOn:   isOn,
		reason: reason,
	}
}

func (d FeatureFlagDecision) IsOn() bool {
	return d.isOn
}

func (d FeatureFlagDecision) Reason() string {
	return d.reason
}

func (d FeatureFlagDecision) String() string {
	return fmt.Sprintf("FeatureFlagDecision(isOn=%t, reason=%s, config=%s)", d.IsOn(), d.Reason(), d.Config)
}

type RemoteConfigDecision struct {
	value  interface{}
	reason string
}

func NewRemoteConfigDecision(value interface{}, reason string) RemoteConfigDecision {
	return RemoteConfigDecision{
		value:  value,
		reason: reason,
	}
}

func (d RemoteConfigDecision) Value() interface{} {
	return d.value
}

func (d RemoteConfigDecision) Reason() string {
	return d.reason
}

func (d RemoteConfigDecision) String() string {
	return fmt.Sprintf("RemoteConfigDecision(value=%s, reason=%s)", d.Value(), d.Reason())
}

const (
	ReasonSdkNotReady                    = "SDK_NOT_READY"
	ReasonException                      = "EXCEPTION"
	ReasonInvalidInput                   = "INVALID_INPUT"
	ReasonExperimentNotFound             = "EXPERIMENT_NOT_FOUND"
	ReasonExperimentDraft                = "EXPERIMENT_DRAFT"
	ReasonExperimentPaused               = "EXPERIMENT_PAUSED"
	ReasonExperimentCompleted            = "EXPERIMENT_COMPLETED"
	ReasonOverridden                     = "OVERRIDDEN"
	ReasonTrafficNotAllocated            = "TRAFFIC_NOT_ALLOCATED"
	ReasonNotInMutualExclusionExperiment = "NOT_IN_MUTUAL_EXCLUSION_EXPERIMENT"
	ReasonIdentifierNotFound             = "IDENTIFIER_NOT_FOUND"
	ReasonVariationDropped               = "VARIATION_DROPPED"
	ReasonTrafficAllocated               = "TRAFFIC_ALLOCATED"
	ReasonTrafficAllocatedByTargeting    = "TRAFFIC_ALLOCATED_BY_TARGETING"
	ReasonNotInExperimentTarget          = "NOT_IN_EXPERIMENT_TARGET"
	ReasonFeatureFlagNotFound            = "FEATURE_FLAG_NOT_FOUND"
	ReasonFeatureFlagInactive            = "FEATURE_FLAG_INACTIVE"
	ReasonIndividualTargetMatch          = "INDIVIDUAL_TARGET_MATCH"
	ReasonTargetRuleMatch                = "TARGET_RULE_MATCH"
	ReasonDefaultRule                    = "DEFAULT_RULE"
	ReasonRemoteConfigParameterNotFound  = "REMOTE_CONFIG_PARAMETER_NOT_FOUND"
	ReasonTypeMismatch                   = "TYPE_MISMATCH"
)
