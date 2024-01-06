package client

import (
	"net"
	"syscall"
)

func ClassifyNetworkError(err error) string {
	cause := err
	for {
		// Unwrap
		if unwrap, ok := cause.(interface{ Unwrap() error }); ok {
			cause = unwrap.Unwrap()
			continue
		}
		break
	}

	// DNSError.IsNotFound
	if cause, ok := cause.(*net.DNSError); ok && cause.Err == "no such host" {
		return "name not found"
	}

	if cause, ok := cause.(syscall.Errno); ok {
		if cause == 10061 || cause == syscall.ECONNREFUSED {
			return "connection refused"
		}
	}

	if cause, ok := cause.(net.Error); ok && cause.Timeout() {
		return "timeout"
	}

	return err.Error()
}
