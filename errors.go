package dnsrewrite

import "errors"

var (
	ErrClosed      = errors.New("dnsrewrite: engine closed")
	ErrNilEngine   = errors.New("dnsrewrite: nil engine")
	ErrNilMatcher  = errors.New("dnsrewrite: nil matcher")
	ErrNilQuestion = errors.New("dnsrewrite: nil question")
	ErrInvalidRule = errors.New("dnsrewrite: invalid rule")
	ErrNotFound    = errors.New("dnsrewrite: rule not found")
	ErrUpstream    = errors.New("dnsrewrite: upstream failure")
	ErrCanceled    = errors.New("dnsrewrite: canceled")
	ErrTimeout     = errors.New("dnsrewrite: timeout")
	ErrPersist     = errors.New("dnsrewrite: persist failure")
	ErrRefuse      = errors.New("dnsrewrite: query refused")
	ErrNoMatch     = errors.New("dnsrewrite: no matching rule")
)
