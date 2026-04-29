package httpx

import (
	"context"
	"net/http"

	"github.com/arcgolabs/collectionx/list"
	"github.com/arcgolabs/collectionx/mapping"
	"github.com/danielgtaylor/huma/v2"
	humaconditional "github.com/danielgtaylor/huma/v2/conditional"
	humasse "github.com/danielgtaylor/huma/v2/sse"
	"github.com/samber/lo"
)

// Docs renderer constants mirror Huma's built-in renderer options.
const (
	DocsRendererScalar            = huma.DocsRendererScalar
	DocsRendererStoplightElements = huma.DocsRendererStoplightElements
	DocsRendererSwaggerUI         = huma.DocsRendererSwaggerUI
)

// HTTP method aliases used by the route registration helpers.
const (
	MethodGet     = http.MethodGet
	MethodPost    = http.MethodPost
	MethodPut     = http.MethodPut
	MethodDelete  = http.MethodDelete
	MethodPatch   = http.MethodPatch
	MethodHead    = http.MethodHead
	MethodOptions = http.MethodOptions
)

// RouteInfo describes a registered route for diagnostics and tests.
type RouteInfo struct {
	Method      string             `json:"method"`
	Path        string             `json:"path"`
	HandlerName string             `json:"handler_name"`
	Comment     string             `json:"comment,omitempty"`
	Tags        *list.List[string] `json:"tags,omitempty"`
}

// String returns related data.
func (r RouteInfo) String() string {
	return r.Method + " " + r.Path + " -> " + r.HandlerName
}

// TypedHandler is the typed handler signature used by `httpx` routes.
type TypedHandler[I, O any] func(ctx context.Context, input *I) (*O, error)

// ConditionalParams aliases Huma conditional request params.
type ConditionalParams = humaconditional.Params

// SSEMessage aliases Huma SSE message for streaming payloads.
type SSEMessage = humasse.Message

// SSESender aliases Huma SSE sender for streaming events.
type SSESender = humasse.Sender

// SSEHandler is the typed handler signature used by SSE routes.
type SSEHandler[I any] func(ctx context.Context, input *I, send SSESender)

// OperationOption mutates a Huma operation before registration.
type OperationOption func(*huma.Operation)

// OpenAPI collection aliases keep public config typed on collectionx.
type (
	OpenAPITags                 = *list.List[string]
	OpenAPITagDefinitions       = *list.List[*huma.Tag]
	OpenAPIParameters           = *list.List[*huma.Param]
	OpenAPIExtensions           = *mapping.Map[string, any]
	OpenAPISecurityScopes       = *list.List[string]
	OpenAPISecurityRequirement  = *mapping.Map[string, OpenAPISecurityScopes]
	OpenAPISecurityRequirements = *list.List[OpenAPISecurityRequirement]
	OpenAPISecuritySchemes      = *mapping.Map[string, *huma.SecurityScheme]
)

// Tags creates an OpenAPI tag list backed by collectionx.
func Tags(values ...string) OpenAPITags {
	return list.NewList(values...)
}

// TagDefinitions creates OpenAPI tag metadata definitions backed by collectionx.
func TagDefinitions(values ...*huma.Tag) OpenAPITagDefinitions {
	return list.NewList(values...)
}

// Parameters creates an OpenAPI parameter list backed by collectionx.
func Parameters(values ...*huma.Param) OpenAPIParameters {
	return list.NewList(values...)
}

// Extensions creates an OpenAPI extension map backed by collectionx.
func Extensions(values map[string]any) OpenAPIExtensions {
	return mapping.NewMapFrom(values)
}

// SecurityScopes creates an OpenAPI scope list backed by collectionx.
func SecurityScopes(values ...string) OpenAPISecurityScopes {
	return list.NewList(values...)
}

// SecurityRequirement creates one OpenAPI security requirement entry.
func SecurityRequirement(name string, scopes ...string) OpenAPISecurityRequirement {
	requirement := mapping.NewMap[string, OpenAPISecurityScopes]()
	if name != "" {
		requirement.Set(name, SecurityScopes(scopes...))
	}
	return requirement
}

// SecurityRequirementMap creates one OpenAPI security requirement from a built-in map.
func SecurityRequirementMap(values map[string][]string) OpenAPISecurityRequirement {
	requirement := mapping.NewMapWithCapacity[string, OpenAPISecurityScopes](len(values))
	lo.ForEach(lo.Entries(values), func(entry lo.Entry[string, []string], _ int) {
		if entry.Key == "" {
			return
		}
		requirement.Set(entry.Key, SecurityScopes(entry.Value...))
	})
	return requirement
}

// SecurityRequirements creates a list of OpenAPI security requirements backed by collectionx.
func SecurityRequirements(values ...OpenAPISecurityRequirement) OpenAPISecurityRequirements {
	return list.NewList(values...)
}

// SecuritySchemes creates an OpenAPI security scheme map backed by collectionx.
func SecuritySchemes(values map[string]*huma.SecurityScheme) OpenAPISecuritySchemes {
	return mapping.NewMapFrom(values)
}

// SecurityOptions configures OpenAPI security schemes and default requirements.
type SecurityOptions struct {
	Schemes      OpenAPISecuritySchemes
	Requirements OpenAPISecurityRequirements
}
