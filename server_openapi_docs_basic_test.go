package httpx_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/arcgolabs/httpx/adapter"
	adapterstd "github.com/arcgolabs/httpx/adapter/std"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
)

func TestServer_WithOpenAPIInfo_UpdatesDocument(t *testing.T) {
	server := newServer(WithOpenAPIInfo("Arc API", "2.0.0", "typed service"))

	openAPI := server.OpenAPI()
	if assert.NotNil(t, openAPI) && assert.NotNil(t, openAPI.Info) {
		assert.Equal(t, "Arc API", openAPI.Info.Title)
		assert.Equal(t, "2.0.0", openAPI.Info.Version)
		assert.Equal(t, "typed service", openAPI.Info.Description)
	}
}

func TestServer_AdapterDocsDisabled_HidesDefaultDocsRoutes(t *testing.T) {
	server := newServer(WithAdapter(adapterstd.New(nil, adapter.HumaOptions{
		DisableDocsRoutes: true,
	})))

	docsReq := newTestRequest(http.MethodGet, "/docs", nil)
	docsRec := serveRequest(t, server, docsReq)
	assert.Equal(t, http.StatusNotFound, docsRec.Code)

	openAPIReq := newTestRequest(http.MethodGet, "/openapi.json", nil)
	openAPIRec := serveRequest(t, server, openAPIReq)
	assert.Equal(t, http.StatusNotFound, openAPIRec.Code)
}

func TestServer_ConfigureOpenAPI_PatchesDocument(t *testing.T) {
	server := newServer()
	server.ConfigureOpenAPI(func(doc *huma.OpenAPI) {
		doc.Tags = append(doc.Tags, &huma.Tag{Name: "internal"})
	})

	openAPI := server.OpenAPI()
	if assert.NotNil(t, openAPI) {
		assert.Len(t, openAPI.Tags, 1)
		assert.Equal(t, "internal", openAPI.Tags[0].Name)
	}
}

func TestGroup_HumaMiddlewareAndModifier(t *testing.T) {
	server := newServer(WithBasePath("/api"))
	group := server.Group("/v1")
	group.UseHumaMiddleware(func(ctx huma.Context, next func(huma.Context)) {
		ctx.AppendHeader("X-Group", "v1")
		next(ctx)
	})
	group.UseSimpleOperationModifier(func(op *huma.Operation) {
		op.Tags = append(op.Tags, "group-tag")
	})

	err := GroupGet(group, "/items", func(_ context.Context, _ *struct{}) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "ok"
		return out, nil
	})
	assert.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/api/v1/items", nil)
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "v1", rec.Header().Get("X-Group"))

	pathItem := server.OpenAPI().Paths["/api/v1/items"]
	if assert.NotNil(t, pathItem) && assert.NotNil(t, pathItem.Get) {
		assert.Contains(t, pathItem.Get.Tags, "group-tag")
	}
}

func TestServer_AdapterDocs_CustomPaths(t *testing.T) {
	server := newServer(WithAdapter(adapterstd.New(nil, adapter.HumaOptions{
		DocsPath:     "/reference",
		OpenAPIPath:  "/spec",
		SchemasPath:  "/contracts",
		DocsRenderer: DocsRendererScalar,
	})))

	docsReq := newTestRequest(http.MethodGet, "/reference", nil)
	docsRec := serveRequest(t, server, docsReq)
	assert.Equal(t, http.StatusOK, docsRec.Code)

	oldDocsReq := newTestRequest(http.MethodGet, "/docs", nil)
	oldDocsRec := serveRequest(t, server, oldDocsReq)
	assert.Equal(t, http.StatusNotFound, oldDocsRec.Code)

	specReq := newTestRequest(http.MethodGet, "/spec.json", nil)
	specRec := serveRequest(t, server, specReq)
	assert.Equal(t, http.StatusOK, specRec.Code)
	assert.Contains(t, specRec.Body.String(), "\"openapi\"")
}
