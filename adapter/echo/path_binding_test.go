package echo_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	httpx "github.com/arcgolabs/httpx"
	echoadapter "github.com/arcgolabs/httpx/adapter/echo"
	"github.com/stretchr/testify/assert"
)

func TestAdapter_StrongTypedPathBinding(t *testing.T) {
	adapter := echoadapter.New(nil)
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

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/users/77", http.NoBody)
	rec := httptest.NewRecorder()
	adapter.Router().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"id":77`)
}

func TestAdapter_CatchAllPathBinding(t *testing.T) {
	adapter := echoadapter.New(nil)
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
	rec := httptest.NewRecorder()
	adapter.Router().ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"bucket":"photos"`)
	assert.Contains(t, rec.Body.String(), `"key":"a/b c"`)
}
