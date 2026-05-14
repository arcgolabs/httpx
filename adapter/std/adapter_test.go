package std_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	stdadapter "github.com/arcgolabs/httpx/adapter/std"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type pingOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func TestAdapter_RouterServesTypedRoute(t *testing.T) {
	a := stdadapter.New(nil)
	huma.Register(a.HumaAPI(), huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
	}, func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "pong"
		return out, nil
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/ping", http.NoBody)
	rec := httptest.NewRecorder()
	a.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "pong")
}

func TestAdapter_CatchAllPathBinding(t *testing.T) {
	a := stdadapter.New(nil)

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

	huma.Register(a.HumaAPI(), huma.Operation{
		OperationID: "get-object",
		Method:      http.MethodGet,
		Path:        "/{bucket}/{key...}",
	}, func(_ context.Context, input *in) (*out, error) {
		result := &out{}
		result.Body.Bucket = input.Bucket
		result.Body.Key = input.Key
		return result, nil
	})

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/photos/a%2Fb%20c", http.NoBody)
	rec := httptest.NewRecorder()
	a.Router().ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"bucket":"photos"`)
	assert.Contains(t, rec.Body.String(), `"key":"a/b c"`)
}
