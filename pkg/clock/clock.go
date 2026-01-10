package clock

import "time"

// Clock abstracts time for deterministic testing
type Clock interface {
	Now() time.Time
}
