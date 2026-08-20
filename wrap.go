package dnsrewrite

import (
	"context"
	"fmt"

	ierr "example.com/dnsrewrite/internal/errors"
)

func wrapUpstream(err error) error {
	if err == nil {
		return nil
	}
	return ierr.WrapErr(ErrUpstream, err)
}

func wrapCanceled(err error) error {
	if err == nil {
		return ErrCanceled
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return fmt.Errorf("%w: %w", ErrCanceled, err)
	}
	return fmt.Errorf("%w: %w", ErrCanceled, err)
}

func wrapPersist(err error) error {
	if err == nil {
		return nil
	}
	return ierr.WrapErr(ErrPersist, err)
}
