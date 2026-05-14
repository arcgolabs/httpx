package httpx_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_StrongTypedPathBindingOnStdAdapter(t *testing.T) {
	server := newServer()

	type in struct {
		UserID int `path:"id"`
	}
	type out struct {
		Body struct {
			ID int `json:"id"`
		}
	}

	err := Get(server, "/users/{id}", func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.ID = input.UserID
		return result, nil
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/users/123", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":123`)
}

func TestServer_CatchAllPathBindingOnStdAdapter(t *testing.T) {
	server := newServer()

	type in struct {
		Bucket string `path:"bucket"`
		Key    string `path:"key"`
	}
	type out struct {
		Body struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		}
	}

	err := Get(server, "/{bucket}/{key...}", func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.Bucket = input.Bucket
		result.Body.Key = input.Key
		return result, nil
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/photos/a%2Fb%20c", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"bucket":"photos"`)
	assert.Contains(t, w.Body.String(), `"key":"a/b c"`)

	route, ok := server.MatchRoute(http.MethodGet, "/photos/a/b/c")
	assert.True(t, ok)
	assert.Equal(t, "/{bucket}/{key...}", route.Path)

	_, ok = server.MatchRoute(http.MethodGet, "/photos")
	assert.False(t, ok)

	route, ok = server.MatchRoute(http.MethodGet, "/photos/")
	assert.True(t, ok)
	assert.Equal(t, "/{bucket}/{key...}", route.Path)

	assert.NotNil(t, server.OpenAPI().Paths["/{bucket}/{key}"])
	assert.Nil(t, server.OpenAPI().Paths["/{bucket}/{key...}"])
}

func TestServer_CatchAllPathTailProvidesRawValue(t *testing.T) {
	server := newServer()

	type in struct {
		Bucket string   `path:"bucket"`
		Key    PathTail `path:"key"`
	}
	type out struct {
		Body struct {
			Bucket  string `json:"bucket"`
			Key     string `json:"key"`
			Raw     string `json:"raw"`
			Present bool   `json:"present"`
		}
	}

	err := Get(server, "/{bucket}/{key...}", func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.Bucket = input.Bucket
		result.Body.Key = input.Key.String()
		result.Body.Raw = input.Key.Raw()
		result.Body.Present = input.Key.Present()
		return result, nil
	})
	require.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/photos/a%2Fb%20c", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"bucket":"photos"`)
	assert.Contains(t, w.Body.String(), `"key":"a/b c"`)
	assert.Contains(t, w.Body.String(), `"raw":"a%2Fb%20c"`)
	assert.Contains(t, w.Body.String(), `"present":true`)

	pathItem := server.OpenAPI().Paths["/{bucket}/{key}"]
	require.NotNil(t, pathItem)
	require.NotNil(t, pathItem.Get)
	param := findOpenAPIPathParameter(pathItem.Get.Parameters, "key")
	require.NotNil(t, param)
	assert.Equal(t, true, param.Extensions["x-httpx-catch-all"])
}

func TestServer_CatchAllPathBindingWithBasePathAndGroup(t *testing.T) {
	server := newServer(WithBasePath("/api"))
	group := server.Group("/s3")

	type in struct {
		Bucket string `path:"bucket"`
		Key    string `path:"key"`
	}
	type out struct {
		Body struct {
			Bucket string `json:"bucket"`
			Key    string `json:"key"`
		}
	}

	err := GroupGet(group, "/{bucket}/{key...}", func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.Bucket = input.Bucket
		result.Body.Key = input.Key
		return result, nil
	})
	require.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/api/s3/photos/a/b", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"bucket":"photos"`)
	assert.Contains(t, w.Body.String(), `"key":"a/b"`)

	route, ok := server.MatchRoute(http.MethodGet, "/api/s3/photos/a/b")
	assert.True(t, ok)
	assert.Equal(t, "/api/s3/{bucket}/{key...}", route.Path)
	assert.NotNil(t, server.OpenAPI().Paths["/api/s3/{bucket}/{key}"])
}

func TestServer_CatchAllRouteDoesNotShadowStaticRoute(t *testing.T) {
	server := newServer()

	type catchAllIn struct {
		Bucket string `path:"bucket"`
		Key    string `path:"key"`
	}
	type out struct {
		Body struct {
			Kind string `json:"kind"`
		}
	}

	err := Get(server, "/{bucket}/{key...}", func(_ context.Context, _ *catchAllIn) (*out, error) {
		result := &out{}
		result.Body.Kind = "catch-all"
		return result, nil
	})
	require.NoError(t, err)

	err = Get(server, "/assets/health", func(_ context.Context, _ *struct{}) (*out, error) {
		result := &out{}
		result.Body.Kind = "static"
		return result, nil
	})
	require.NoError(t, err)

	route, ok := server.MatchRoute(http.MethodGet, "/assets/health")
	assert.True(t, ok)
	assert.Equal(t, "/assets/health", route.Path)

	req := newTestRequest(http.MethodGet, "/assets/health", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"kind":"static"`)
}

func TestServer_CatchAllPathMustBeLast(t *testing.T) {
	server := newServer()

	type in struct {
		Key string `path:"key"`
		ID  string `path:"id"`
	}

	err := Get(server, "/files/{key...}/{id}", func(_ context.Context, _ *in) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "ok"
		return out, nil
	})

	assert.Error(t, err)
}

func findOpenAPIPathParameter(params []*huma.Param, name string) *huma.Param {
	for _, param := range params {
		if param != nil && param.In == "path" && param.Name == name {
			return param
		}
	}
	return nil
}
