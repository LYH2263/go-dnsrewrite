package clock

import "time"

// Real 系统时钟。
type Real struct{}

func (Real) Now() time.Time                  { return time.Now() }
func (Real) Since(t time.Time) time.Duration { return time.Since(t) }
