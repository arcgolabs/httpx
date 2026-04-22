package httpx_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
)

type fakeFiberAdapterNoApp struct{}

func (f *fakeFiberAdapterNoApp) Name() string { return "fiber" }

func (f *fakeFiberAdapterNoApp) HumaAPI() huma.API { return nil }

func (f *fakeFiberAdapterNoApp) Listen(_ string) error { return ErrAdapterNotFound }

type fakeAdapterWithoutHuma struct{}

func (f *fakeAdapterWithoutHuma) Name() string { return "fake" }

func (f *fakeAdapterWithoutHuma) HumaAPI() huma.API { return nil }

type fakeLifecycleAdapter struct {
	listenAddr string
	shutdown   bool
}

func (f *fakeLifecycleAdapter) Name() string { return "lifecycle" }

func (f *fakeLifecycleAdapter) HumaAPI() huma.API { return nil }

func (f *fakeLifecycleAdapter) Listen(addr string) error {
	f.listenAddr = addr
	return nil
}

func (f *fakeLifecycleAdapter) Shutdown() error {
	f.shutdown = true
	return nil
}

func TestServer_GenericHandlerReturnsHTTPXError(t *testing.T) {
	server := newServer()
	err := Get(server, "/forbidden", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		return nil, NewError(http.StatusForbidden, "no permission")
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/forbidden", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "no permission")
}

func TestServer_GenericHandlerPanic(t *testing.T) {
	server := newServer()
	err := Get(server, "/panic", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		panic("boom")
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/panic", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, strings.ToLower(w.Body.String()), "panic in handler")
}

func TestServer_GenericHandlerNilOutputReturnsNoContent(t *testing.T) {
	server := newServer()
	err := Get(server, "/empty", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		return &pingOutput{}, nil
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/empty", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestServer_AdapterWithoutHumaAPI(t *testing.T) {
	server := newServer(
		WithAdapter(&fakeAdapterWithoutHuma{}),
	)

	err := Get(server, "/huma", func(_ context.Context, _ *struct{}) (*humaPingOutput, error) {
		out := &humaPingOutput{}
		out.Body.Message = "ok"
		return out, nil
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, ErrAdapterNotFound.Error())
}

func TestServer_ListenAndServe_FiberWithoutApp(t *testing.T) {
	server := newServer(WithAdapter(&fakeFiberAdapterNoApp{}))
	err := server.ListenAndServe(":0")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAdapterNotFound))
}

func TestServer_ListenPort_UsesPortShortcut(t *testing.T) {
	lifecycle := &fakeLifecycleAdapter{}
	server := newServer(WithAdapter(lifecycle))

	err := server.ListenPort(8080)
	assert.NoError(t, err)
	assert.Equal(t, ":8080", lifecycle.listenAddr)
}

func TestServer_ListenPort_InvalidPort(t *testing.T) {
	server := newServer()

	err := server.ListenPort(-1)
	assert.EqualError(t, err, fmt.Sprintf("httpx: invalid port %d", -1))
}

func TestServer_Listen_DelegatesToAdapter(t *testing.T) {
	lifecycle := &fakeLifecycleAdapter{}
	server := newServer(WithAdapter(lifecycle))

	err := server.Listen(":9090")
	assert.NoError(t, err)
	assert.Equal(t, ":9090", lifecycle.listenAddr)
}

func TestServer_Shutdown_DelegatesToAdapter(t *testing.T) {
	lifecycle := &fakeLifecycleAdapter{}
	server := newServer(WithAdapter(lifecycle))

	err := server.Shutdown()
	assert.NoError(t, err)
	assert.True(t, lifecycle.shutdown)
}
