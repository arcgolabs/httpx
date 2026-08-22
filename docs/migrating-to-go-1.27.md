---
title: 'Migrating to Go 1.27'
linkTitle: 'migrating-to-go-1.27'
description: 'Breaking changes introduced by the Go 1.27 and dependency upgrade'
weight: 8
---

## Migrating to Go 1.27

This release requires Go 1.27.0 and updates all direct dependencies to their latest compatible releases at the time of the migration.

### Typed route methods

`New` now returns `*httpx.Server`, making Go 1.27 generic methods available directly on servers and groups:

```go
server := httpx.New(httpx.WithAdapter(runtime))
server.MustGet("/users/{id}", getUser)

v1 := server.Group("/v1")
v1.MustPost("/users", createUser)
```

The existing package-level `Get`, `Post`, `GroupGet`, and related helpers remain available.

### Echo v5

The Echo adapter now uses `github.com/labstack/echo/v5`. Update application imports and middleware to the v5 module path. The adapter serves the Echo engine through `net/http`; normal shutdown and context cancellation now return `nil`.

### ETagSet

`ETagSet` is no longer a slice alias. Use `ParseETagSet` to construct it and `Values` to obtain a defensive copy of its tags:

```go
tags, err := httpx.ParseETagSet(`"current", W/"previous"`)
if err != nil {
	return err
}
for _, tag := range tags.Values() {
	// ...
}
```

This representation prevents Huma's slice element decoding from interpreting a complete conditional-request header as one tag per list element.

### Repository workspace

Internal modules do not declare other modules from this repository in their `go.mod` files. Repository builds and tests must run with the checked-in `go.work`, which owns internal module resolution.
