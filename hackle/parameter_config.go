package hackle

type ParameterConfig interface {
	GetString(key string, defaultValue string) string
	GetNumber(key string, defaultValue float64) float64
	GetBool(key string, defaultValue bool) bool
}
