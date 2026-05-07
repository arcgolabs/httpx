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
