package clock

import "time"

type Clock interface {
	CurrentMillis() int64
	Tick() int64
}

var System = SystemClock{}

type SystemClock struct{}

func (c SystemClock) CurrentMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (c SystemClock) Tick() int64 {
	return time.Now().UnixNano()
}

type FixedClock struct {
	time int64
}

func Fixed(time int) Clock {
	return &FixedClock{int64(time)}
}

func (c *FixedClock) CurrentMillis() int64 {
	return c.time
}

func (c *FixedClock) Tick() int64 {
	return c.time
}
