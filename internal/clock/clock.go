package clock

import "time"

// Clock 可注入时钟。
type Clock interface {
	Now() time.Time
	Since(t time.Time) time.Duration
}
