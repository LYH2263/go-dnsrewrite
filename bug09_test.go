package dnsrewrite_test

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"testing"

	"example.com/dnsrewrite"
	"example.com/dnsrewrite/internal/upstream"
)

type trackConn struct {
	net.Conn
	closed bool
}

func (t *trackConn) Close() error {
	t.closed = true
	return t.Conn.Close()
}

type trackDialer struct {
	last *trackConn
	real upstream.Dialer
}

func (d *trackDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	c, err := d.real.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	tc := &trackConn{Conn: c}
	d.last = tc
	return tc, nil
}

type plainDialer struct {
	d net.Dialer
}

func (p *plainDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return p.d.DialContext(ctx, network, address)
}

func TestBug09_UpstreamConnClosed(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		var hdr [4]byte
		if _, err := io.ReadFull(conn, hdr[:]); err != nil {
			return
		}
		n := binary.BigEndian.Uint32(hdr[:])
		body := make([]byte, n)
		_, _ = io.ReadFull(conn, body)
		ans, _ := json.Marshal(upstream.Answer{ID: 1, RCode: 0})
		out := make([]byte, 4+len(ans))
		binary.BigEndian.PutUint32(out[:4], uint32(len(ans)))
		copy(out[4:], ans)
		_, _ = conn.Write(out)
	}()

	addr := ln.Addr().String()
	e := dnsrewrite.New(dnsrewrite.WithDefaultUpstream(addr))
	defer e.Close()
	td := &trackDialer{real: &plainDialer{d: net.Dialer{}}}
	e.Upstream().SetDialer(td)

	_, err = e.Resolve(dnsrewrite.Question{Name: "a.local", Type: dnsrewrite.TypeA})
	if err != nil {
		t.Fatal(err)
	}
	if td.last == nil || !td.last.closed {
		t.Fatal("upstream connection was not Closed after Engine.Resolve")
	}
}
