package clock

import "time"

// FakeClock is a controllable clock for tests
type FakeClock struct {
	current time.Time
}

func NewFakeClock(start time.Time) *FakeClock {
	return &FakeClock{
		current: start,
	}
}

func (f *FakeClock) Now() time.Time {
	return f.current
}

// Advance moves the clock forward by duration
func (f *FakeClock) Advance(d time.Duration) {
	f.current = f.current.Add(d)
}
