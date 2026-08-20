package upstream

import (
	"context"
	"net"
	"time"
)

// Dialer 可注入拨号。
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type netDialer struct {
	d net.Dialer
}

func (n *netDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return n.d.DialContext(ctx, network, address)
}

func defaultDialer(timeout time.Duration) Dialer {
	return &netDialer{d: net.Dialer{Timeout: timeout}}
}
