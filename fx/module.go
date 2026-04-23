package fx

import (
	pkgfx "github.com/arcgolabs/pkg/fx"
	"github.com/arcgolabs/httpx"
	"go.uber.org/fx"
)

// ServerParams defines parameters for httpx fxx module.
type ServerParams struct {
	fx.In

	// Options grouped from WithServerOptions and NewHttpxModule arguments.
	Options []httpx.ServerOption `group:"httpx_server_options"`
}

// ServerResult defines result for httpx fxx module.
type ServerResult struct {
	fx.Out

	// Server is the created httpx server runtime.
	Server httpx.ServerRuntime
}

// NewServer creates a httpx server runtime from grouped options.
func NewServer(params ServerParams) ServerResult {
	return ServerResult{Server: httpx.New(params.Options...)}
}

// WithServerOptions adds server options into fxx option group.
func WithServerOptions(opts ...httpx.ServerOption) fx.Option {
	return pkgfx.ProvideOptionGroup[httpx.Server, httpx.ServerOption]("httpx_server_options", opts...)
}

// NewHttpxModule creates a httpx fxx module.
// It reuses httpx.ServerOption as the module input options.
func NewHttpxModule(opts ...httpx.ServerOption) fx.Option {
	return fx.Module("httpx",
		fx.Provide(NewServer),
		WithServerOptions(opts...),
	)
}
