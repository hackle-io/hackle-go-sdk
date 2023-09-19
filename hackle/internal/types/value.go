package types

type ValueType string

func (t ValueType) String() string {
	return string(t)
}

const (
	String  ValueType = "STRING"
	Number  ValueType = "NUMBER"
	Bool    ValueType = "BOOLEAN"
	Version ValueType = "VERSION"
	Json    ValueType = "JSON"
)

var types = map[string]ValueType{
	string(String):  String,
	string(Number):  Number,
	string(Bool):    Bool,
	string(Version): Version,
	string(Json):    Json,
}

func TypeFrom(value string) (ValueType, bool) {
	valueType, ok := types[value]
	return valueType, ok
}
