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
	// 已取消的 ctx 立即返回，避免死等到超时才醒。
	if err := ctx.Err(); err != nil {
		return ierr.WrapErr(ierr.ErrCanceled, err)
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	select {
	case <-ready:
		return nil
	case <-timer.C:
		return ierr.ErrTimeout
	case <-ctx.Done():
		return ierr.WrapErr(ierr.ErrCanceled, ctx.Err())
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
	// 已取消的 ctx 立即返回，不进入轮询。
	if err := ctx.Err(); err != nil {
		return ierr.WrapErr(ierr.ErrCanceled, err)
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		// 轮询间隔内也要能被 ctx.Done 打断，而非死等到下一次 tick。
		select {
		case <-ctx.Done():
			return ierr.WrapErr(ierr.ErrCanceled, ctx.Err())
		case <-t.C:
			if ready != nil && ready() {
				return nil
			}
		}
	}
}
