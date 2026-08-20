package errors

import "errors"

var (
	ErrCanceled = errors.New("dnsrewrite/internal: canceled")
	ErrTimeout  = errors.New("dnsrewrite/internal: timeout")
	ErrUpstream = errors.New("dnsrewrite/internal: upstream")
	ErrPersist  = errors.New("dnsrewrite/internal: persist")
)
