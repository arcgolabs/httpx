## Overview

`httpx` is a lightweight HTTP service organization layer built on top of Huma.
It gives you a stable **server/group/endpoint** API surface across multiple runtimes (std/chi, gin, echo, fiber), while still allowing direct access to Huma when you need it.

## Install

```bash
go get github.com/arcgolabs/httpx@latest
```

## Current capabilities

- Unified typed route registration across adapters (`Get`, `Post`, `Put`, `Patch`, `Delete`...)
- Adapter-based runtime integration (`std`, `gin`, `echo`, `fiber`)
- Endpoint/module organization with scoped registrars, metadata defaults, and limited auto route sugar (`Endpoint`, `Registrar`, `EndpointSpec`, `Auto`)
- First-class OpenAPI and documentation control (docs route exposure is adapter-owned)
- Typed Server-Sent Events (SSE) (`GetSSE`, `GroupGetSSE`)
- Policy-based route capabilities (`RouteWithPolicies`, `GroupRouteWithPolicies`)
- Conditional request handling (`If-Match`, `If-None-Match`, `If-Modified-Since`, `If-Unmodified-Since`)
- Route-aware observability helpers for Prometheus and OpenTelemetry
- `dix` integration helpers for lifecycle and listen wiring
- Direct Huma escape hatches (`HumaAPI`, `OpenAPI`, `ConfigureOpenAPI`)
- Optional request validation via `go-playground/validator`
- Route introspection API for testing and diagnostics

## Package layout

- Core: `github.com/arcgolabs/httpx`
- Adapters:
    - `github.com/arcgolabs/httpx/adapter/std`
    - `github.com/arcgolabs/httpx/adapter/gin`
    - `github.com/arcgolabs/httpx/adapter/echo`
    - `github.com/arcgolabs/httpx/adapter/fiber` (Fiber v3)
- Optional:
    - `github.com/arcgolabs/httpx/middleware`
    - `github.com/arcgolabs/httpx/dix`
    - `github.com/arcgolabs/httpx/websocket`

## Documentation map (recommended reading)

- Minimal typed server: [Getting Started](./getting-started)
- Adapter wiring: [Adapters](./adapters)
- Endpoint organization: [Endpoint Organization](./endpoint-organization)
- Middleware and observability: [Middleware and Observability](./middleware-and-observability)
- DI wiring: [dix Integration](./dix-integration)
- OpenAPI and docs: [OpenAPI and docs](./openapi-and-docs)

## Runnable examples (repository)

- Quickstart: [examples/quickstart](https://github.com/arcgolabs/httpx/tree/main/examples/quickstart)
- Adapters:
    - [examples/std](https://github.com/arcgolabs/httpx/tree/main/examples/std)
    - [examples/gin](https://github.com/arcgolabs/httpx/tree/main/examples/gin)
    - [examples/echo](https://github.com/arcgolabs/httpx/tree/main/examples/echo)
    - [examples/fiber](https://github.com/arcgolabs/httpx/tree/main/examples/fiber)
- Auth / organization:
    - [examples/auth](https://github.com/arcgolabs/httpx/tree/main/examples/auth)
    - [examples/organization](https://github.com/arcgolabs/httpx/tree/main/examples/organization)
- Streaming:
    - SSE: [examples/sse](https://github.com/arcgolabs/httpx/tree/main/examples/sse)
    - Websocket: [examples/websocket](https://github.com/arcgolabs/httpx/tree/main/examples/websocket)
- Conditional requests: [examples/conditional](https://github.com/arcgolabs/httpx/tree/main/examples/conditional)
- Endpoint registration: [examples/endpoint](https://github.com/arcgolabs/httpx/tree/main/examples/endpoint)
- `dix` backend wiring: [github.com/arcgolabs/dix](https://github.com/arcgolabs/dix)

## Positioning (how to think about it)

- `Huma`: typed operations, schemas, OpenAPI/docs, middleware model
- `adapter/*`: runtime/router integration + native middleware ecosystem
- `httpx`: unified service organization API + exposes selected Huma capabilities
- `httpx/websocket`: lightweight websocket helper used directly at framework/router layer; it is intentionally not part of the typed route API
