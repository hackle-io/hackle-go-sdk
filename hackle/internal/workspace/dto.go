package workspace

//goland:noinspection GoNameStartsWithPackageName
type WorkspaceDTO struct {
	Experiments             []ExperimentDTO             `json:"experiments"`
	FeatureFlags            []ExperimentDTO             `json:"featureFlags"`
	Buckets                 []BucketDTO                 `json:"buckets"`
	Events                  []EventTypeDTO              `json:"events"`
	Segments                []SegmentDTO                `json:"segments"`
	Containers              []ContainerDTO              `json:"containers"`
	ParameterConfigurations []ParameterConfigurationDTO `json:"parameterConfigurations"`
	RemoteConfigParameters  []RemoteConfigParameterDTO  `json:"remoteConfigParameters"`
}

type ExperimentDTO struct {
	ID                int64          `json:"id"`
	Key               int64          `json:"key"`
	Name              *string        `json:"name"`
	Status            string         `json:"status"`
	Version           int            `json:"version"`
	Variations        []VariationDTO `json:"variations"`
	Execution         ExecutionDTO   `json:"execution"`
	WinnerVariationID *int64         `json:"winnerVariationId"`
	IdentifierType    string         `json:"identifierType"`
	ContainerID       *int64         `json:"containerId"`
}

type VariationDTO struct {
	ID                       int64  `json:"id"`
	Key                      string `json:"key"`
	Status                   string `json:"status"`
	ParameterConfigurationID *int64 `json:"parameterConfigurationId"`
}

type ExecutionDTO struct {
	Status           string            `json:"status"`
	Version          int               `json:"version"`
	UserOverrides    []UserOverrideDTO `json:"userOverrides"`
	SegmentOverrides []TargetRuleDTO   `json:"segmentOverrides"`
	TargetAudiences  []TargetDTO       `json:"targetAudiences"`
	TargetRules      []TargetRuleDTO   `json:"targetRules"`
	DefaultRule      TargetActionDTO   `json:"defaultRule"`
}

type UserOverrideDTO struct {
	UserID      string `json:"userId"`
	VariationID int64  `json:"variationId"`
}

type BucketDTO struct {
	ID       int64     `json:"id"`
	Seed     int       `json:"seed"`
	SlotSize int       `json:"slotSize"`
	Slots    []SlotDTO `json:"slots"`
}

type SlotDTO struct {
	StartInclusive int   `json:"startInclusive"`
	EndExclusive   int   `json:"endExclusive"`
	VariationID    int64 `json:"variationId"`
}

type EventTypeDTO struct {
	ID  int64  `json:"id"`
	Key string `json:"key"`
}

type TargetDTO struct {
	Conditions []TargetConditionDTO `json:"conditions"`
}

type TargetConditionDTO struct {
	Key   TargetKeyDTO   `json:"key"`
	Match TargetMatchDTO `json:"match"`
}

type TargetKeyDTO struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type TargetMatchDTO struct {
	Type      string        `json:"type"`
	Operator  string        `json:"operator"`
	ValueType string        `json:"valueType"`
	Values    []interface{} `json:"values"`
}

type TargetActionDTO struct {
	Type        string `json:"type"`
	VariationID *int64 `json:"variationId"`
	BucketID    *int64 `json:"bucketId"`
}

type TargetRuleDTO struct {
	Target TargetDTO       `json:"target"`
	Action TargetActionDTO `json:"action"`
}

type SegmentDTO struct {
	ID      int64       `json:"id"`
	Key     string      `json:"key"`
	Type    string      `json:"type"`
	Targets []TargetDTO `json:"targets"`
}

type ContainerDTO struct {
	ID            int64               `json:"id"`
	EnvironmentID int64               `json:"environmentId"`
	BucketID      int64               `json:"bucketId"`
	Groups        []ContainerGroupDTO `json:"groups"`
}

type ContainerGroupDTO struct {
	ID          int64   `json:"id"`
	Experiments []int64 `json:"experiments"`
}

type ParameterConfigurationDTO struct {
	ID         int64          `json:"id"`
	Parameters []ParameterDTO `json:"parameters"`
}

type ParameterDTO struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

type RemoteConfigParameterDTO struct {
	ID             int64                       `json:"id"`
	Key            string                      `json:"key"`
	Type           string                      `json:"type"`
	IdentifierType string                      `json:"identifierType"`
	TargetRules    []RemoteConfigTargetRuleDTO `json:"targetRules"`
	DefaultValue   RemoteConfigValueDTO        `json:"defaultValue"`
}

type RemoteConfigTargetRuleDTO struct {
	Key      string               `json:"key"`
	Name     string               `json:"name"`
	Target   TargetDTO            `json:"target"`
	BucketID int64                `json:"bucketId"`
	Value    RemoteConfigValueDTO `json:"value"`
}

type RemoteConfigValueDTO struct {
	ID    int64       `json:"id"`
	Value interface{} `json:"value"`
}
