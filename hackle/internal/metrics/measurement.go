package metrics

type Measurement struct {
	Field     Field
	valueFunc func() float64
}

func NewMeasurement(field Field, valueFunc func() float64) Measurement {
	return Measurement{
		Field:     field,
		valueFunc: valueFunc,
	}
}

func (m *Measurement) Value() float64 {
	return m.valueFunc()
}

type Field string

const (
	FieldCount Field = "count"
	FieldTotal Field = "total"
	FieldMax   Field = "max"
	FieldMean  Field = "mean"
)

func (f Field) String() string {
	return string(f)
}
