package model

type Sdk struct {
	Key     string
	Name    string
	Version string
}

func NewSdk(sdkKey string) Sdk {
	return Sdk{
		Key:     sdkKey,
		Name:    SdkName,
		Version: SdkVersion,
	}
}

const (
	SdkName    = "go-sdk"
	SdkVersion = "0.0.1-SNAPSHOT"
)
