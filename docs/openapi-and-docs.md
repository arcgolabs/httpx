---
title: 'httpx OpenAPI and Docs'
linkTitle: 'openapi-and-docs'
description: 'Configure OpenAPI metadata and docs routes via adapter options'
weight: 4
---

## OpenAPI and docs

With Huma under the hood, `httpx` exposes OpenAPI and documentation routes via adapter configuration (for example `adapter.HumaOptions`).

In `std` adapter options you typically set:

- `DocsPath` (for example `/docs`)
- `OpenAPIPath` (for example `/openapi.json`)
- `Title`, `Version`, `Description`

## Minimal example

```go
package main

import (
	"context"
	"net/http"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/go-chi/chi/v5"
)

type out struct {
	Body struct {
		Status string `json:"status"`
	} `json:"body"`
}

func main() {
	router := chi.NewMux()

	stdAdapter := std.New(router, adapter.HumaOptions{
		Title:       "My API",
		Version:     "1.0.0",
		Description: "Service API",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	s := httpx.New(httpx.WithAdapter(stdAdapter))
	httpx.MustGet(s, "/health", func(ctx context.Context, _ *struct{}) (*out, error) {
		o := &out{}
		o.Body.Status = "ok"
		return o, nil
	})

	_ = http.ListenAndServe(":8080", router)
}
```

## Related

- [Getting Started](./getting-started)
- OpenAPI-heavy runnable samples: [examples/httpx/endpoint](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/endpoint), [examples/httpx/organization](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/organization)

