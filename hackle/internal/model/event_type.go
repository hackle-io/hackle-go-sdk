package model

type EventType struct {
	ID  int64
	Key string
}

func NewUndefinedEvent(key string) EventType {
	return EventType{ID: 0, Key: key}
}
