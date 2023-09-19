package mocks

func CreateEvent(key string) Event {
	return Event{key: key}
}

type Event struct {
	key        string
	value      float64
	properties map[string]interface{}
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
