package httpx_test

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func TestServer_WithPanicRecover_Enabled(t *testing.T) {
	server := newServer(WithPanicRecover(true))

	err := Get(server, "/panic", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		panic("boom")
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/panic", nil)
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "panic in handler: boom")
}

func TestServer_WithPanicRecover_Disabled(t *testing.T) {
	server := newServer(WithPanicRecover(false))

	err := Get(server, "/panic", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		panic("boom")
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/panic", nil)

	assert.Panics(t, func() {
		_ = serveRequest(t, server, req)
	})
}

func TestServer_WithAccessLog_LogsRequests(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	server := newServer(
		WithLogger(logger),
		WithAccessLog(true),
	)

	type in struct {
		ID int `path:"id"`
	}

	err := Get(server, "/users/{id}", func(_ context.Context, _ *in) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "ok"
		return out, nil
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/users/42", nil)
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	output := logs.String()
	assert.True(t, strings.Contains(output, "\"msg\":\"httpx request\""))
	assert.True(t, strings.Contains(output, "\"method\":\"GET\""))
	assert.True(t, strings.Contains(output, "\"path\":\"/users/42\""))
	assert.True(t, strings.Contains(output, "\"status\":200"))
	assert.True(t, strings.Contains(output, "\"route\":\"/users/{id}\""))
}

func TestServer_WithPrintRoutes_LogsOnRegistration(t *testing.T) {
	var logs bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	server := newServer(
		WithLogger(logger),
		WithPrintRoutes(true),
	)

	err := Get(server, "/routes", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "ok"
		return out, nil
	})
	assert.NoError(t, err)

	output := logs.String()
	assert.Contains(t, output, "Registered routes")
	assert.Contains(t, output, "GET /routes")
}
