package clock

import "time"

type Clock interface {
	Now() time.Time
}

type MockClock interface {
	Now() time.Time
	Set(t time.Time)
}

type RealClock struct{}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

type FakeClock struct {
	t time.Time
}

func (c *FakeClock) Now() time.Time {
	return c.t
}

func (c *FakeClock) Set(t time.Time) {
	c.t = t
}
