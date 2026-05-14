package adapter

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

const catchAllMetadataKey = "httpx.catch_all_path_param"

// CatchAllOpenAPIExtension marks an OpenAPI path parameter as an httpx
// catch-all tail parameter.
const CatchAllOpenAPIExtension = "x-httpx-catch-all"

// RouterPathKind identifies the router syntax behind a Huma adapter.
type RouterPathKind string

const (
	RouterPathChi   RouterPathKind = "chi"
	RouterPathGin   RouterPathKind = "gin"
	RouterPathEcho  RouterPathKind = "echo"
	RouterPathFiber RouterPathKind = "fiber"
)

// CatchAllParam describes a logical catch-all path parameter.
type CatchAllParam struct {
	Name       string
	NativeName string
	Index      int
}

// CompiledRoutePath is the adapter-specific route path and parameter mapping.
type CompiledRoutePath struct {
	OperationPath string
	RouterPath    string
	CatchAll      *CatchAllParam
}

// PathAdapterOption configures the path bridge adapter.
type PathAdapterOption func(*PathAdapter)

// PathAdapter wraps a Huma adapter and translates httpx path extensions before
// they reach the underlying router.
type PathAdapter struct {
	base                   huma.Adapter
	kind                   RouterPathKind
	missingCatchAllHandler func(huma.Context) bool
}

var _ huma.Adapter = (*PathAdapter)(nil)

// NewPathAdapter wraps a Huma adapter with httpx path translation.
func NewPathAdapter(base huma.Adapter, kind RouterPathKind, opts ...PathAdapterOption) *PathAdapter {
	adapter := &PathAdapter{
		base: base,
		kind: kind,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(adapter)
		}
	}
	return adapter
}

// WithMissingCatchAllHandler handles router matches where the catch-all segment
// was not actually present. Fiber needs this because `/*` also matches the
// slashless parent route.
func WithMissingCatchAllHandler(handler func(huma.Context) bool) PathAdapterOption {
	return func(adapter *PathAdapter) {
		adapter.missingCatchAllHandler = handler
	}
}

// RouterOnlyConfig disables Huma's generated docs/schema routes for an internal
// adapter API used only to obtain its router bridge.
func RouterOnlyConfig(cfg huma.Config) huma.Config {
	cfg.DocsPath = ""
	cfg.OpenAPIPath = ""
	cfg.SchemasPath = ""
	return cfg
}

// Handle implements huma.Adapter.
func (a *PathAdapter) Handle(op *huma.Operation, handler func(ctx huma.Context)) {
	if a == nil || a.base == nil {
		return
	}

	compiled := CompileRouterOperationPath(op, a.kind)
	routerOp := *op
	routerOp.Path = compiled.RouterPath
	a.base.Handle(&routerOp, func(ctx huma.Context) {
		if compiled.CatchAll != nil {
			if !CatchAllRemainderPresent(ctx, compiled) &&
				a.missingCatchAllHandler != nil &&
				a.missingCatchAllHandler(ctx) {
				return
			}
			ctx = WrapCatchAllContext(ctx, op, compiled)
		}
		handler(ctx)
	})
}

// ServeHTTP implements huma.Adapter.
func (a *PathAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a == nil || a.base == nil {
		http.NotFound(w, r)
		return
	}
	a.base.ServeHTTP(w, r)
}

// NormalizeOperationPath rewrites `{name...}` to OpenAPI-compatible `{name}`
// while recording enough metadata for the adapter bridge to register a router
// wildcard and still bind `path:"name"`.
func NormalizeOperationPath(op *huma.Operation) error {
	if op == nil {
		return nil
	}

	normalized, name, err := normalizeCatchAllPath(op.Path)
	if err != nil {
		return err
	}
	if name == "" {
		return nil
	}

	op.Path = normalized
	if op.Metadata == nil {
		op.Metadata = map[string]any{}
	}
	op.Metadata[catchAllMetadataKey] = name
	return nil
}

// MarkCatchAllPathParameter annotates the documented OpenAPI parameter that
// represents an httpx catch-all path tail.
func MarkCatchAllPathParameter(op *huma.Operation) {
	name := catchAllNameFromOperation(op)
	if name == "" {
		return
	}
	for _, param := range op.Parameters {
		if param == nil || param.In != "path" || param.Name != name {
			continue
		}
		if param.Extensions == nil {
			param.Extensions = map[string]any{}
		}
		param.Extensions[CatchAllOpenAPIExtension] = true
		return
	}
}

// CompileRouterOperationPath converts a Huma/OpenAPI path into the syntax
// expected by the underlying Huma router adapter.
func CompileRouterOperationPath(op *huma.Operation, kind RouterPathKind) CompiledRoutePath {
	if op == nil {
		return CompiledRoutePath{}
	}

	result := CompiledRoutePath{
		OperationPath: op.Path,
		RouterPath:    op.Path,
	}

	name := catchAllNameFromOperation(op)
	path := op.Path
	if name == "" {
		normalized, parsedName, err := normalizeCatchAllPath(path)
		if err != nil || parsedName == "" {
			return result
		}
		name = parsedName
		path = normalized
	}

	routerPath, index, ok := replaceCatchAllSegment(path, name, catchAllRouterSegment(kind, name))
	if !ok {
		return result
	}

	result.OperationPath = path
	result.RouterPath = routerPath
	result.CatchAll = &CatchAllParam{
		Name:       name,
		NativeName: catchAllNativeName(kind, name),
		Index:      index,
	}
	return result
}

func catchAllName(metadata map[string]any) string {
	if metadata == nil {
		return ""
	}
	name, ok := metadata[catchAllMetadataKey].(string)
	if !ok {
		return ""
	}
	return name
}

func catchAllRouterSegment(kind RouterPathKind, name string) string {
	switch kind {
	case RouterPathChi, RouterPathEcho, RouterPathFiber:
		return "*"
	case RouterPathGin:
		return "*" + name
	default:
		return "*"
	}
}

func catchAllNativeName(kind RouterPathKind, name string) string {
	switch kind {
	case RouterPathChi, RouterPathEcho, RouterPathFiber:
		return "*"
	case RouterPathGin:
		return name
	default:
		return "*"
	}
}
