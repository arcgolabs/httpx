package httpx

import (
	"github.com/arcgolabs/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
)

type catchAllDocumentAPI struct {
	huma.API
}

func registerAPIForOperation(api huma.API, op huma.Operation) huma.API {
	if op.Path == "" {
		return api
	}
	return catchAllDocumentAPI{API: api}
}

func (api catchAllDocumentAPI) DocumentOperation(op *huma.Operation) {
	adapter.MarkCatchAllPathParameter(op)
	if documenter, ok := api.API.(huma.OperationDocumenter); ok {
		documenter.DocumentOperation(op)
		return
	}
	if !op.Hidden {
		api.OpenAPI().AddOperation(op)
	}
}
