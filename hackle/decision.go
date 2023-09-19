package hackle

import "fmt"

type ExperimentDecision interface {
	fmt.Stringer
	ParameterConfig
	Variation() string
	Reason() string
}

type FeatureFlagDecision interface {
	fmt.Stringer
	ParameterConfig
	IsOn() bool
	Reason() string
}

type RemoteConfigDecision interface {
	fmt.Stringer
	Value() interface{}
	Reason() string
}
