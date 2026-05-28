package fiber_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpx "github.com/arcgolabs/httpx"
	fiberadapter "github.com/arcgolabs/httpx/adapter/fiber"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdapter_StrongTypedPathBinding(t *testing.T) {
	adapter := fiberadapter.New(nil)
	server := httpx.New(httpx.WithAdapter(adapter))

	type in struct {
		UserID int `path:"id"`
	}
	type out struct {
		Body struct {
			ID int `json:"id"`
		}
	}

	err := httpx.Get(server, "/users/{id}", func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.ID = input.UserID
		return result, nil
	})
	assert.NoError(t, err)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users/66", http.NoBody)
	resp, err := testRequest(adapter, req)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), `"id":66`)
}

func TestAdapter_CatchAllPathBinding(t *testing.T) {
	adapter := fiberadapter.New(nil)
	server := httpx.New(httpx.WithAdapter(adapter))

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

	err := httpx.Get(server, "/{bucket}/{key...}", func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.Bucket = input.Bucket
		result.Body.Key = input.Key
		return result, nil
	})
	assert.NoError(t, err)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/photos/a%2Fb%20c", http.NoBody)
	resp, err := testRequest(adapter, req)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), `"bucket":"photos"`)
	assert.Contains(t, string(body), `"key":"a/b c"`)
}

func TestAdapter_CatchAllDoesNotCaptureBareParentPath(t *testing.T) {
	adapter := fiberadapter.New(nil)
	server := httpx.New(httpx.WithAdapter(adapter))

	type objectIn struct {
		Bucket string `path:"bucket"`
		Key    string `path:"key"`
	}
	type bucketIn struct {
		Bucket string `path:"bucket"`
	}
	type out struct {
		Body struct {
			Kind   string `json:"kind"`
			Bucket string `json:"bucket"`
			Key    string `json:"key,omitempty"`
		}
	}

	err := httpx.Get(server, "/{bucket}/{key...}", func(_ context.Context, input *objectIn) (*out, error) {
		result := &out{}
		result.Body.Kind = "object"
		result.Body.Bucket = input.Bucket
		result.Body.Key = input.Key
		return result, nil
	})
	require.NoError(t, err)

	err = httpx.Get(server, "/{bucket}", func(_ context.Context, input *bucketIn) (*out, error) {
		result := &out{}
		result.Body.Kind = "bucket"
		result.Body.Bucket = input.Bucket
		return result, nil
	})
	require.NoError(t, err)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/photos", http.NoBody)
	resp, err := testRequest(adapter, req)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), `"kind":"bucket"`)
	assert.Contains(t, string(body), `"bucket":"photos"`)
	assert.NotContains(t, string(body), `"kind":"object"`)
}
