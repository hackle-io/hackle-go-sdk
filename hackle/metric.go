package hackle

import (
	"github.com/hackle-io/hackle-go-sdk/hackle/internal/metrics"
	"strconv"
)

func recordExperiment(sample *metrics.TimerSample, key int64, decision ExperimentDecision) {
	tags := metrics.Tags{
		"key":       strconv.FormatInt(key, 10),
		"variation": decision.Variation(),
		"reason":    decision.Reason(),
	}
	timer := metrics.NewTimer("experiment.decision", tags)
	sample.Stop(timer)
}

func recordFeatureFlag(sample *metrics.TimerSample, key int64, decision FeatureFlagDecision) {
	tags := metrics.Tags{
		"key":    strconv.FormatInt(key, 10),
		"on":     strconv.FormatBool(decision.IsOn()),
		"reason": decision.Reason(),
	}
	timer := metrics.NewTimer("feature.flag.decision", tags)
	sample.Stop(timer)
}

func recordRemoteConfig(sample *metrics.TimerSample, key string, decision RemoteConfigDecision) {
	tags := metrics.Tags{
		"key":    key,
		"reason": decision.Reason(),
	}
	timer := metrics.NewTimer("remote.config.decision", tags)
	sample.Stop(timer)
}
