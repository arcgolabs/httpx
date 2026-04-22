package shared

import (
	"log/slog"

	"github.com/arcgolabs/httpx"
	"github.com/arcgolabs/httpx/adapter"
)

// NewRuntime builds a shared demo server runtime around the provided adapter.
func NewRuntime(a adapter.Host, logger *slog.Logger) httpx.ServerRuntime {
	return httpx.New(
		httpx.WithAdapter(a),
		httpx.WithLogger(logger),
		httpx.WithPrintRoutes(true),
		httpx.WithValidation(),
	)
}
