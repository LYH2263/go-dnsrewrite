package dnsrewrite_test

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"testing"
	"time"

	"example.com/dnsrewrite"
)

func TestBug07_ResolveContextHonorsCancel(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	started := make(chan struct{})
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		close(started)
		// 慢读：等调用方取消
		buf := make([]byte, 4)
		_, _ = io.ReadFull(conn, buf)
		n := binary.BigEndian.Uint32(buf)
		body := make([]byte, n)
		_, _ = io.ReadFull(conn, body)
		time.Sleep(3 * time.Second)
		ans, _ := json.Marshal(map[string]any{"ID": 1, "RCode": 0, "Answers": []any{}})
		hdr := make([]byte, 4)
		binary.BigEndian.PutUint32(hdr, uint32(len(ans)))
		_, _ = conn.Write(append(hdr, ans...))
	}()

	e := dnsrewrite.New(dnsrewrite.WithUpstreamTimeout(5 * time.Second))
	defer e.Close()
	_ = e.AddRule(dnsrewrite.RuleSpec{
		ID: "slow", Pattern: "slow.local", Kind: dnsrewrite.MatchExact,
		Action: dnsrewrite.ActionForward, Upstream: ln.Addr().String(), Enabled: true,
	})

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := e.ResolveContext(ctx, dnsrewrite.Question{Name: "slow.local", Type: dnsrewrite.TypeA})
		done <- err
	}()
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("upstream not hit")
	}
	cancel()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("ResolveContext ignored cancel")
	}
}
