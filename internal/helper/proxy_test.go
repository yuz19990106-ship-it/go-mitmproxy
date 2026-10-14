package helper

import (
	"context"
	"net/url"
	"testing"
	"time"
)

// Regression test for upstream issue lqqyt2423/go-mitmproxy#110:
// GetProxyConn panicked with a nil pointer dereference when the SOCKS5
// proxy URL dial failed, because the error path called conn.Close() on
// a nil conn. It must return the dial error instead of panicking.
func TestGetProxyConnSocks5DialFailureReturnsError(t *testing.T) {
	proxyUrl, err := url.Parse("socks5://127.0.0.1:1")
	if err != nil {
		t.Fatalf("parse proxy url: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := GetProxyConn(ctx, proxyUrl, "example.com:443", false)
	if err == nil {
		if conn != nil {
			conn.Close()
		}
		t.Fatal("expected dial error, got nil")
	}
	if conn != nil {
		conn.Close()
		t.Fatal("expected nil conn on dial error")
	}
	t.Logf("dial failed as expected: %v", err)
}
