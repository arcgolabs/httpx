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
	resp, err := adapter.Router().Test(req, -1)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, resp.Body.Close())
	}()

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(body), `"id":66`)
}
