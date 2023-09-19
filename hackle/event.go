package hackle

import "github.com/hackle-io/hackle-go-sdk/hackle/internal/properties"

type Event struct {
	key        string
	value      float64
	properties map[string]interface{}
}

func NewEvent(key string) Event {
	return NewEventBuilder(key).Build()
}

func (e Event) Key() string {
	return e.key
}

func (e Event) Value() float64 {
	return e.value
}

func (e Event) Properties() map[string]interface{} {
	return e.properties
}

type EventBuilder struct {
	key        string
	value      float64
	properties *properties.Builder
}

func NewEventBuilder(key string) *EventBuilder {
	return &EventBuilder{
		key:        key,
		properties: properties.NewBuilder(),
	}
}

func (b *EventBuilder) Value(value float64) *EventBuilder {
	b.value = value
	return b
}

func (b *EventBuilder) Property(key string, value interface{}) *EventBuilder {
	b.properties.Add(key, value)
	return b
}

func (b *EventBuilder) Properties(properties map[string]interface{}) *EventBuilder {
	b.properties.AddAll(properties)
	return b
}

func (b *EventBuilder) Build() Event {
	return Event{
		key:        b.key,
		value:      b.value,
		properties: b.properties.Build(),
	}
}
