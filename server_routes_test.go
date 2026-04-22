package httpx_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_MatchRoute_ExactRouteWinsOverParameterizedRoute(t *testing.T) {
	server := newServer()

	require.NoError(t, Get(server, "/users/{id}", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "param"
		return out, nil
	}))
	require.NoError(t, Get(server, "/users/me", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "exact"
		return out, nil
	}))

	matched, ok := matchRoute(server, http.MethodGet, "/users/me")
	require.True(t, ok)
	assert.Equal(t, "/users/me", matched.Path)

	matched, ok = matchRoute(server, http.MethodGet, "/users/42")
	require.True(t, ok)
	assert.Equal(t, "/users/{id}", matched.Path)
}

func TestServer_MatchRoute_OverlappingParameterizedRoutesKeepRegistrationOrder(t *testing.T) {
	server := newServer()

	require.NoError(t, Get(server, "/{kind}/list", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "generic"
		return out, nil
	}))
	require.NoError(t, Get(server, "/users/{id}", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "specific"
		return out, nil
	}))

	matched, ok := matchRoute(server, http.MethodGet, "/users/list")
	require.True(t, ok)
	assert.Equal(t, "/{kind}/list", matched.Path)
}

func TestServer_MatchRoute_OverlappingParameterizedRoutesKeepRegistrationOrderWhenReversed(t *testing.T) {
	server := newServer()

	require.NoError(t, Get(server, "/users/{id}", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "specific"
		return out, nil
	}))
	require.NoError(t, Get(server, "/{kind}/list", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "generic"
		return out, nil
	}))

	matched, ok := matchRoute(server, http.MethodGet, "/users/list")
	require.True(t, ok)
	assert.Equal(t, "/users/{id}", matched.Path)
}

func TestServer_AddTag_ReplacesExistingTagByName(t *testing.T) {
	server := newServer()

	server.AddTag(&huma.Tag{Name: "users", Description: "first"})
	server.AddTag(&huma.Tag{Name: "users", Description: "updated"})

	doc := server.OpenAPI()
	require.NotNil(t, doc)
	require.Len(t, doc.Tags, 1)
	assert.Equal(t, "updated", doc.Tags[0].Description)
}

func TestServer_DuplicateRouteRegistrationReturnsError(t *testing.T) {
	server := newServer()

	require.NoError(t, Get(server, "/users", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "first"
		return out, nil
	}))

	err := Get(server, "/users", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "second"
		return out, nil
	})
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrRouteAlreadyExists)
	assert.Equal(t, 1, server.RouteCount())

	req := newTestRequest(http.MethodGet, "/users", nil)
	rec := serveRequest(t, server, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "first")
}
