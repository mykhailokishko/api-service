package domain

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}

func NewRealClock() RealClock {
	return RealClock{}
}

type MockClock struct {
	FixedTime time.Time
}

func NewMockClock(t time.Time) MockClock {
	return MockClock{
		FixedTime: t,
	}
}

func (mc MockClock) Now() time.Time {
	return mc.FixedTime
}
