package httpx_test

import (
	"context"
	"testing"

	"github.com/arcgolabs/httpx"
	"github.com/stretchr/testify/require"
)

func TestServer_GenericRouteMethods(t *testing.T) {
	server := httpx.New()

	err := server.Get("/method", func(_ context.Context, _ *struct{}) (*struct{}, error) {
		return &struct{}{}, nil
	})
	require.NoError(t, err)
	require.True(t, server.HasRoute(httpx.MethodGet, "/method"))
}

func TestGroup_GenericRouteMethods(t *testing.T) {
	server := httpx.New()
	group := server.Group("/v1")

	err := group.Post("/method", func(_ context.Context, _ *struct{}) (*struct{}, error) {
		return &struct{}{}, nil
	})
	require.NoError(t, err)
	require.True(t, server.HasRoute(httpx.MethodPost, "/v1/method"))
}
