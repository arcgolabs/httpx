---
title: 'httpx Getting Started'
linkTitle: 'getting-started'
description: 'Create a typed server, register routes, enable validation, and serve docs'
weight: 2
---

## Getting Started

`httpx` is a lightweight HTTP service organization layer built on top of Huma. You pick an adapter (`std`, `gin`, `echo`, `fiber`) and register typed routes via `httpx.Get/Post/...` (or `MustGet/MustPost/...`).

This page shows a minimal server with:

- `std` adapter (chi router)
- typed endpoints (`/health`, `/v1/users`, `/v1/users/{id}`)
- OpenAPI + docs routes
- request validation (`WithValidation()`)

## 1) Install

```bash
go get github.com/DaiYuANg/arcgo/httpx@latest
go get github.com/go-chi/chi/v5
```

## 2) Create `main.go`

```go
package main

import (
	"context"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/go-chi/chi/v5"
)

type healthOutput struct {
	Body struct {
		Status string `json:"status"`
	} `json:"body"`
}

type createUserInput struct {
	Body struct {
		Name  string `json:"name" validate:"required,min=2,max=64"`
		Email string `json:"email" validate:"required,email"`
	} `json:"body"`
}

type createUserOutput struct {
	Body struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"body"`
}

type getUserInput struct {
	ID int `path:"id"`
}

type getUserOutput struct {
	Body struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"body"`
}

func main() {
	router := chi.NewMux()

	stdAdapter := std.New(router, adapter.HumaOptions{
		Title:       "httpx getting-started",
		Version:     "1.0.0",
		Description: "Typed HTTP example",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	server := httpx.New(
		httpx.WithAdapter(stdAdapter),
		httpx.WithBasePath("/api"),
		httpx.WithValidation(),
	)

	httpx.MustGet(server, "/health", func(ctx context.Context, _ *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})

	v1 := server.Group("/v1")

	httpx.MustGroupPost(v1, "/users", func(ctx context.Context, in *createUserInput) (*createUserOutput, error) {
		out := &createUserOutput{}
		out.Body.ID = 1001
		out.Body.Name = in.Body.Name
		out.Body.Email = in.Body.Email
		return out, nil
	})

	httpx.MustGroupGet(v1, "/users/{id}", func(ctx context.Context, in *getUserInput) (*getUserOutput, error) {
		out := &getUserOutput{}
		out.Body.ID = in.ID
		out.Body.Name = "demo-user"
		return out, nil
	})

	_ = server.ListenPort(8080)
}
```

## Next

- Adapter choices and wiring: [Adapters](./adapters)
- OpenAPI and docs control: [OpenAPI and docs](./openapi-and-docs)

