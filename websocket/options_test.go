// Package websocket_test contains tests for websocket.
package websocket_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/arcgolabs/httpx/websocket"
)

func TestDefaultOptions(t *testing.T) {
	opts := websocket.DefaultOptions()
	if opts.HandshakeTimeout <= 0 {
		t.Fatalf("expected handshake timeout > 0, got %v", opts.HandshakeTimeout)
	}
	if opts.MaxMessageSize <= 0 {
		t.Fatalf("expected max message size > 0, got %d", opts.MaxMessageSize)
	}
}

func TestApplyOptions(t *testing.T) {
	checkOrigin := func(*http.Request) bool { return true }
	opts := websocket.DefaultOptions()
	for _, option := range []websocket.Option{
		websocket.WithHandshakeTimeout(2 * time.Second),
		websocket.WithReadTimeout(3 * time.Second),
		websocket.WithWriteTimeout(4 * time.Second),
		websocket.WithIdleTimeout(5 * time.Second),
		websocket.WithMaxMessageSize(128),
		websocket.WithCompression(true),
		websocket.WithCheckOrigin(checkOrigin),
	} {
		option(&opts)
	}

	if opts.HandshakeTimeout != 2*time.Second {
		t.Fatalf("unexpected handshake timeout: %v", opts.HandshakeTimeout)
	}
	if opts.ReadTimeout != 3*time.Second {
		t.Fatalf("unexpected read timeout: %v", opts.ReadTimeout)
	}
	if opts.WriteTimeout != 4*time.Second {
		t.Fatalf("unexpected write timeout: %v", opts.WriteTimeout)
	}
	if opts.IdleTimeout != 5*time.Second {
		t.Fatalf("unexpected idle timeout: %v", opts.IdleTimeout)
	}
	if opts.MaxMessageSize != 128 {
		t.Fatalf("unexpected max message size: %d", opts.MaxMessageSize)
	}
	if !opts.EnableCompression {
		t.Fatal("expected compression to be enabled")
	}
	if opts.CheckOrigin == nil || !opts.CheckOrigin(nil) {
		t.Fatal("expected check origin function to be set")
	}
}

func TestApplyOptionsFallback(t *testing.T) {
	opts := websocket.DefaultOptions()
	websocket.WithHandshakeTimeout(0)(&opts)
	websocket.WithMaxMessageSize(0)(&opts)
	if opts.HandshakeTimeout <= 0 {
		t.Fatalf("expected fallback handshake timeout, got %v", opts.HandshakeTimeout)
	}
	if opts.MaxMessageSize <= 0 {
		t.Fatalf("expected fallback max message size, got %d", opts.MaxMessageSize)
	}
}
