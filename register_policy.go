package httpx

import (
	"strings"

	"github.com/samber/lo"
	"github.com/samber/oops"
)

// RouteWithPolicies registers a typed handler with runtime/OpenAPI policies.
func RouteWithPolicies[I, O any](s ServerRuntime, method, path string, handler TypedHandler[I, O], policies ...RoutePolicy[I, O]) error {
	server := unwrapServer(s)
	if server == nil {
		return oops.In("httpx").
			With("op", "route_with_policies", "method", strings.ToUpper(method), "path", path, "policy_count", len(policies)).
			Wrapf(ErrRouteNotRegistered, "validate server")
	}
	fullPath := joinRoutePath(server.basePath, path)
	return registerTyped(server, server.HumaAPI(), method, fullPath, fullPath, handler, nil, policies)
}

// GroupRouteWithPolicies registers a grouped typed handler with policies.
func GroupRouteWithPolicies[I, O any](g *Group, method, path string, handler TypedHandler[I, O], policies ...RoutePolicy[I, O]) error {
	if g == nil || g.server == nil {
		return oops.In("httpx").
			With("op", "group_route_with_policies", "method", strings.ToUpper(method), "path", path, "policy_count", len(policies)).
			Wrapf(ErrRouteNotRegistered, "validate route group")
	}
	fullPath := joinRoutePath(g.server.basePath, joinRoutePath(g.prefix, path))

	target := g.server.HumaAPI()
	registerPath := fullPath
	if g.humaGroup != nil {
		target = g.humaGroup
		registerPath = path
	}

	return registerTyped(g.server, target, method, registerPath, fullPath, handler, nil, policies)
}

// MustRouteWithPolicies registers a route with policies and panics on failure.
func MustRouteWithPolicies[I, O any](s ServerRuntime, method, path string, handler TypedHandler[I, O], policies ...RoutePolicy[I, O]) {
	lo.Must0(RouteWithPolicies(s, method, path, handler, policies...))
}

// MustGroupRouteWithPolicies registers a grouped route with policies and panics on failure.
func MustGroupRouteWithPolicies[I, O any](g *Group, method, path string, handler TypedHandler[I, O], policies ...RoutePolicy[I, O]) {
	lo.Must0(GroupRouteWithPolicies(g, method, path, handler, policies...))
}
