package upstream

import (
	"context"
	"time"
)

// WaitHealthy 等待地址可拨通（用 WaitReady 语义包装）。
func WaitHealthy(ctx context.Context, dial Dialer, addr string, timeout time.Duration) error {
	ready := make(chan struct{})
	go func() {
		defer close(ready)
		dctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		conn, err := dial.DialContext(dctx, "tcp", addr)
		if err == nil {
			_ = conn.Close()
		}
	}()
	return WaitReady(ctx, ready, timeout)
}
