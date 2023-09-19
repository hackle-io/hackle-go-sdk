package model

type Variation struct {
	ID                       int64
	Key                      string
	IsDropped                bool
	ParameterConfigurationID *int64
}
