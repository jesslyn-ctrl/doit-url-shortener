package clock

import "time"

// SystemClock uses the real system time
type SystemClock struct{}

func NewSystemClock() *SystemClock {
	return &SystemClock{}
}

func (SystemClock) Now() time.Time {
	return time.Now()
}
