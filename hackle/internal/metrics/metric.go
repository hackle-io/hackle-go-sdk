package metrics

import (
	"sort"
	"strings"
)

type Metric interface {
	ID() ID
	Measure() []Measurement
}

type BaseMetric struct {
	id ID
}

func NewBaseMetric(id ID) *BaseMetric {
	return &BaseMetric{id: id}
}

func (m *BaseMetric) ID() ID {
	return m.id
}

type Type string

const (
	TypeCounter Type = "COUNTER"
	TypeTimer   Type = "TIMER"
)

type Tags map[string]string

type ID struct {
	fqn  string
	Name string
	Tags Tags
	Type Type
}

func NewID(name string, tags Tags, Type Type) ID {
	return ID{
		fqn:  fqn(name, tags),
		Name: name,
		Tags: tags,
		Type: Type,
	}
}

func fqn(name string, tags Tags) string {
	if len(tags) == 0 {
		return name + "{}"
	}

	keys := make([]string, 0, len(tags))
	for key := range tags {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	i := make([]string, 0, len(tags))
	for _, key := range keys {
		i = append(i, key+"="+tags[key])
	}

	return name + "{" + strings.Join(i, ", ") + "}"
}
