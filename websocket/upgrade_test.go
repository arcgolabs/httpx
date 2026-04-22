// Package websocket_test contains tests for websocket.
package websocket_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arcgolabs/httpx/websocket"
)

func TestUpgradeNilHandler(t *testing.T) {
	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/ws", http.NoBody)
	rec := httptest.NewRecorder()

	err := websocket.Upgrade(rec, req, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, websocket.ErrUpgradeFailed) {
		t.Fatalf("expected ErrUpgradeFailed, got %v", err)
	}
}
