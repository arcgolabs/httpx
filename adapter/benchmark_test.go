package adapter_test

import (
	"testing"

	adapter "github.com/arcgolabs/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
)

func BenchmarkApplyHumaConfig(b *testing.B) {
	opts := adapter.HumaOptions{
		Description: "benchmark docs",
		DocsPath:    "/docs/v1",
		OpenAPIPath: "/openapi/v1",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for range b.N {
		cfg := huma.DefaultConfig("benchmark", "1.0.0")
		adapter.ApplyHumaConfig(&cfg, opts)
		if cfg.DocsPath == "" || cfg.OpenAPIPath == "" {
			b.Fatal("expected docs/openapi paths to be set")
		}
	}
}
