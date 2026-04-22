package httpx_test

import (
	"context"
	"net/http"
	"testing"

	adapterecho "github.com/arcgolabs/httpx/adapter/echo"
	adapterfiber "github.com/arcgolabs/httpx/adapter/fiber"
	adaptergin "github.com/arcgolabs/httpx/adapter/gin"
	"github.com/stretchr/testify/assert"
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

func TestServer_StrongTypedPathBindingOnGinAdapter(t *testing.T) {
	server := newServer(WithAdapter(adaptergin.New(nil)))

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

	req := newTestRequest(http.MethodGet, "/users/88", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":88`)
}

func TestServer_StrongTypedPathBindingOnEchoAdapter(t *testing.T) {
	server := newServer(WithAdapter(adapterecho.New(nil)))

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

	req := newTestRequest(http.MethodGet, "/users/77", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":77`)
}

func TestServer_StrongTypedPathBindingOnFiberAdapter(t *testing.T) {
	server := newServer(WithAdapter(adapterfiber.New(nil)))

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

	req := newTestRequest(http.MethodGet, "/users/66", nil)
	w := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":66`)
}
