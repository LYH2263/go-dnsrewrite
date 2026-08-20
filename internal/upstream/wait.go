package upstream

import (
	"context"
	"time"

	ierr "example.com/dnsrewrite/internal/errors"
)

// WaitReady 等待 ready 或超时；必须响应 ctx 取消。
func WaitReady(ctx context.Context, ready <-chan struct{}, timeout time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	// BUG: 忽略已取消 ctx，也不在 select 中听 Done
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-ready:
		return nil
	case <-timer.C:
		return ierr.ErrTimeout
	}
}

// WaitPoll 轮询直到 ready 函数返回 true。
func WaitPoll(ctx context.Context, interval time.Duration, ready func() bool) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if ready != nil && ready() {
		return nil
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		// BUG: 轮询不查 ctx
		<-t.C
		if ready != nil && ready() {
			return nil
		}
	}
}
