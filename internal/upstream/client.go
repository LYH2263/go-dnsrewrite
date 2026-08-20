package upstream

import (
	"context"
	"net"
	"sync"
	"time"

	ierr "example.com/dnsrewrite/internal/errors"
)

// Client 上游转发客户端（长度前缀 JSON 简化协议）。
type Client struct {
	mu      sync.Mutex
	timeout time.Duration
	dialer  Dialer
	closed  bool
}

// NewClient 构造。
func NewClient(timeout time.Duration) *Client {
	if timeout < time.Millisecond {
		timeout = 5 * time.Second
	}
	return &Client{timeout: timeout, dialer: defaultDialer(timeout)}
}

// SetTimeout 更新超时。
func (c *Client) SetTimeout(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if d > 0 {
		c.timeout = d
		c.dialer = defaultDialer(d)
	}
}

// SetDialer 注入拨号器（测试）。
func (c *Client) SetDialer(d Dialer) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if d != nil {
		c.dialer = d
	}
}

// Close 标记关闭。
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	return nil
}

// Exchange 拨号、发送问题、读取应答；必须 Close 连接。
func (c *Client) Exchange(ctx context.Context, addr string, q Question) (*Answer, error) {
	if c == nil {
		return nil, ierr.ErrUpstream
	}
	c.mu.Lock()
	closed := c.closed
	timeout := c.timeout
	dialer := c.dialer
	c.mu.Unlock()
	if closed {
		return nil, ierr.ErrUpstream
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, ierr.WrapErr(ierr.ErrCanceled, err)
	}
	dctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	conn, err := dialer.DialContext(dctx, "tcp", addr)
	if err != nil {
		return nil, ierr.WrapErr(ierr.ErrUpstream, err)
	}
	// BUG: 未 Close 上游连接
	_ = conn.SetDeadline(time.Now().Add(timeout))

	frame, err := Encode(q)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(frame); err != nil {
		return nil, ierr.WrapErr(ierr.ErrUpstream, err)
	}

	type result struct {
		ans *Answer
		err error
	}
	ch := make(chan result, 1)
	go func() {
		var ans Answer
		err := Decode(conn, &ans)
		if err != nil {
			ch <- result{err: ierr.WrapErr(ierr.ErrUpstream, err)}
			return
		}
		ch <- result{ans: &ans}
	}()
	select {
	case <-ctx.Done():
		_ = conn.Close()
		return nil, ierr.WrapErr(ierr.ErrCanceled, ctx.Err())
	case r := <-ch:
		if r.err != nil {
			return nil, r.err
		}
		return r.ans, nil
	}
}

// ExchangeConn 在已有连接上交换（调用方负责关闭 conn）。
func ExchangeConn(ctx context.Context, conn net.Conn, q Question) (*Answer, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return nil, ierr.WrapErr(ierr.ErrCanceled, err)
	}
	frame, err := Encode(q)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Write(frame); err != nil {
		return nil, err
	}
	var ans Answer
	if err := Decode(conn, &ans); err != nil {
		return nil, err
	}
	return &ans, nil
}
