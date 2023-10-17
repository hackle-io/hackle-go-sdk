package metrics

type mockRegistry struct {
	*BaseRegistry
	isClosed   bool
	newCounter func(id ID) Counter
	newTimer   func(id ID) Timer
}

func newMockRegistry() *mockRegistry {
	registry := &mockRegistry{}
	registry.BaseRegistry = NewBaseRegistry(registry)
	return registry
}

func (m *mockRegistry) NewCounter(id ID) Counter {
	if m.newCounter != nil {
		return m.newCounter(id)
	}
	return NewNoopCounter(id)
}

func (m *mockRegistry) NewTimer(id ID) Timer {
	if m.newTimer != nil {
		return m.newTimer(id)
	}
	return NewNoopTimer(id)
}

func (m *mockRegistry) Close() {
	m.isClosed = true
}

func (m *mockRegistry) IsClosed() bool {
	return m.isClosed
}
