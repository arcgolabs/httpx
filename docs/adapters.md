---
title: 'httpx Adapters'
linkTitle: 'adapters'
description: 'Choose std/gin/echo/fiber adapters and wire httpx to your router'
weight: 3
---

## Adapters

Adapters integrate Huma + `httpx` with a runtime router/framework.

Available adapters:

- `httpx/adapter/std` (chi + net/http)
- `httpx/adapter/gin`
- `httpx/adapter/echo`
- `httpx/adapter/fiber` (Fiber v3)

You build an adapter, pass it to `httpx.New(httpx.WithAdapter(...))`, and then register routes on the returned server/group.

## Minimal: std adapter with chi middleware

```go
package main

import (
	"context"
	"net/http"

	"github.com/arcgolabs/httpx"
	"github.com/arcgolabs/httpx/adapter"
	"github.com/arcgolabs/httpx/adapter/std"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type out struct {
	Body struct {
		OK bool `json:"ok"`
	} `json:"body"`
}

func main() {
	router := chi.NewMux()
	router.Use(middleware.Logger, middleware.Recoverer, middleware.RequestID)

	stdAdapter := std.New(router, adapter.HumaOptions{
		Title:       "std adapter",
		Version:     "1.0.0",
		Description: "std adapter example",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	s := httpx.New(httpx.WithAdapter(stdAdapter))
	httpx.MustGet(s, "/ping", func(ctx context.Context, _ *struct{}) (*out, error) {
		o := &out{}
		o.Body.OK = true
		return o, nil
	})

	_ = http.ListenAndServe(":8080", router)
}
```

## Runnable adapter examples (repository)

- [examples/std](https://github.com/arcgolabs/httpx/tree/main/examples/std)
- [examples/gin](https://github.com/arcgolabs/httpx/tree/main/examples/gin)
- [examples/echo](https://github.com/arcgolabs/httpx/tree/main/examples/echo)
- [examples/fiber](https://github.com/arcgolabs/httpx/tree/main/examples/fiber)
