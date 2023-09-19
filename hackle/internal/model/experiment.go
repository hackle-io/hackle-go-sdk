package model

type Experiment struct {
	ID                int64
	Key               int64
	Name              *string
	Type              ExperimentType
	IdentifierType    string
	Status            ExperimentStatus
	Version           int
	ExecutionVersion  int
	Variations        []Variation
	UserOverrides     map[string]int64
	SegmentOverrides  []TargetRule
	TargetAudiences   []Target
	TargetRules       []TargetRule
	DefaultRule       Action
	ContainerID       *int64
	WinnerVariationID *int64
}

type ExperimentType string
type ExperimentStatus string

const (
	ExperimentTypeAbTest      ExperimentType = "AB_TEST"
	ExperimentTypeFeatureFlag ExperimentType = "FEATURE_FLAG"
)

const (
	ExperimentStatusDraft     ExperimentStatus = "DRAFT"
	ExperimentStatusRunning   ExperimentStatus = "RUNNING"
	ExperimentStatusPaused    ExperimentStatus = "PAUSED"
	ExperimentStatusCompleted ExperimentStatus = "COMPLETED"
)

func (e *Experiment) WinnerVariation() (Variation, bool) {
	if e.WinnerVariationID != nil {
		return e.GetVariationByID(*e.WinnerVariationID)
	}
	return Variation{}, false
}

func (e *Experiment) GetVariationByID(id int64) (Variation, bool) {
	for _, variation := range e.Variations {
		if variation.ID == id {
			return variation, true
		}
	}
	return Variation{}, false
}

func (e *Experiment) GetVariationByKey(key string) (Variation, bool) {
	for _, variation := range e.Variations {
		if variation.Key == key {
			return variation, true
		}
	}
	return Variation{}, false
}

func NewExperimentStatusFrom(executionStatusCode string) (ExperimentStatus, bool) {
	switch executionStatusCode {
	case "READY":
		return ExperimentStatusDraft, true
	case "RUNNING":
		return ExperimentStatusRunning, true
	case "PAUSED":
		return ExperimentStatusPaused, true
	case "STOPPED":
		return ExperimentStatusCompleted, true
	}
	return "", false
}
