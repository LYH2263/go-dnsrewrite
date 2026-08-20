package upstream

import (
	"context"
	"errors"
	"testing"
	"time"

	ierr "example.com/dnsrewrite/internal/errors"
)

func TestBug08_UpstreamWaitHonorsContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ready := make(chan struct{})
	err := WaitReady(ctx, ready, 2*time.Second)
	if err == nil {
		t.Fatal("WaitReady ignored canceled ctx")
	}
	if !errors.Is(err, ierr.ErrCanceled) && !errors.Is(err, context.Canceled) {
		t.Fatalf("want canceled, got %v", err)
	}
}
