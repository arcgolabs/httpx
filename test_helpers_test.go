//nolint:testpackage // Exposes internal-only test helpers to httpx_test without expanding package API.
package httpx

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	adapterstd "github.com/arcgolabs/httpx/adapter/std"
	"github.com/stretchr/testify/require"
)

func newTestRequest(method, target string, body io.Reader) *http.Request {
	if body == nil {
		body = http.NoBody
	}
	return httptest.NewRequestWithContext(context.Background(), method, target, body)
}

func serveRequest(tb testing.TB, server ServerRuntime, req *http.Request) *httptest.ResponseRecorder {
	tb.Helper()

	impl := unwrapServer(server)
	require.NotNil(tb, impl)
	require.NotNil(tb, impl.adapter)

	switch host := impl.adapter.(type) {
	case *adapterstd.Adapter:
		rec := httptest.NewRecorder()
		host.Router().ServeHTTP(rec, req)
		return rec
	default:
		tb.Fatalf("unsupported adapter type %T", host)
		return nil
	}
}

// NewServerForTest creates a server with internal defaults for external tests.
func NewServerForTest(opts ...ServerOption) *Server {
	return newServer(opts...)
}

// NewRequestForTest builds a request with a background context for external tests.
func NewRequestForTest(method, target string, body io.Reader) *http.Request {
	return newTestRequest(method, target, body)
}

// ServeRequestForTest executes req against server using the concrete adapter router.
func ServeRequestForTest(tb testing.TB, server ServerRuntime, req *http.Request) *httptest.ResponseRecorder {
	tb.Helper()
	return serveRequest(tb, server, req)
}

// MatchRouteForTest exposes route matching for external tests and benchmarks.
func MatchRouteForTest(server *Server, method, path string) (RouteInfo, bool) {
	if server == nil {
		return RouteInfo{}, false
	}
	return server.matchRoute(method, path)
}

// NormalizeRoutePrefixForTest exposes route prefix normalization for external tests.
func NormalizeRoutePrefixForTest(prefix string) string {
	return normalizeRoutePrefix(prefix)
}

// JoinRoutePathForTest exposes route path joining for external tests.
func JoinRoutePathForTest(basePath, path string) string {
	return joinRoutePath(basePath, path)
}

// FreezeServerForTest freezes further configuration mutation for external tests.
func FreezeServerForTest(server *Server) {
	if server != nil {
		server.freezeConfiguration(context.TODO())
	}
}

// UseHostCapabilityForTest exposes adapter capability dispatch for external tests.
func UseHostCapabilityForTest[T any](server ServerRuntime, use func(T)) bool {
	return useHostCapability(server, use)
}
