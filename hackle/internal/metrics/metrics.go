package metrics

var globalRegistry = NewDelegatingRegistry()

func GlobalRegistry() Registry {
	return globalRegistry
}

func AddRegistry(registry Registry) {
	globalRegistry.Add(registry)
}

func NewCounter(name string, tags Tags) Counter {
	return globalRegistry.Counter(name, tags)
}

func NewTimer(name string, tags Tags) Timer {
	return globalRegistry.Timer(name, tags)
}
