package clock

import "time"

// Fake 可推进时钟。
type Fake struct {
	T time.Time
}

func (f *Fake) Now() time.Time {
	if f.T.IsZero() {
		return time.Unix(0, 0).UTC()
	}
	return f.T
}

func (f *Fake) Since(t time.Time) time.Duration { return f.Now().Sub(t) }

func (f *Fake) Advance(d time.Duration) { f.T = f.Now().Add(d) }
