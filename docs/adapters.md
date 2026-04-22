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
- `httpx/adapter/fiber`

You build an adapter, pass it to `httpx.New(httpx.WithAdapter(...))`, and then register routes on the returned server/group.

## Minimal: std adapter with chi middleware

```go
package main

import (
	"context"
	"net/http"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
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

- [examples/httpx/std](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/std)
- [examples/httpx/gin](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/gin)
- [examples/httpx/echo](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/echo)
- [examples/httpx/fiber](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/fiber)

