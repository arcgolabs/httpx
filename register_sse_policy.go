package httpx

import (
	"strings"

	"github.com/samber/lo"
	"github.com/samber/oops"
)

// RouteSSEWithPolicies registers a typed SSE handler with runtime/OpenAPI SSE policies.
func RouteSSEWithPolicies[I any](s ServerRuntime, method, path string, eventTypeMap map[string]any, handler SSEHandler[I], policies ...SSERoutePolicy[I]) error {
	server := unwrapServer(s)
	if server == nil {
		return oops.In("httpx").
			With("op", "route_sse_with_policies", "method", strings.ToUpper(method), "path", path, "policy_count", len(policies), "event_type_count", len(eventTypeMap)).
			Wrapf(ErrRouteNotRegistered, "validate server")
	}
	fullPath := joinRoutePath(server.basePath, path)
	return registerSSE(server, server.HumaAPI(), method, fullPath, fullPath, eventTypeMap, handler, nil, policies)
}

// GroupRouteSSEWithPolicies registers a grouped typed SSE handler with policies.
func GroupRouteSSEWithPolicies[I any](g *Group, method, path string, eventTypeMap map[string]any, handler SSEHandler[I], policies ...SSERoutePolicy[I]) error {
	if g == nil || g.server == nil {
		return oops.In("httpx").
			With("op", "group_route_sse_with_policies", "method", strings.ToUpper(method), "path", path, "policy_count", len(policies), "event_type_count", len(eventTypeMap)).
			Wrapf(ErrRouteNotRegistered, "validate route group")
	}
	fullPath := joinRoutePath(g.server.basePath, joinRoutePath(g.prefix, path))

	target := g.server.HumaAPI()
	registerPath := fullPath
	if g.humaGroup != nil {
		target = g.humaGroup
		registerPath = path
	}

	return registerSSE(g.server, target, method, registerPath, fullPath, eventTypeMap, handler, nil, policies)
}

// MustRouteSSEWithPolicies registers an SSE route with policies and panics on failure.
func MustRouteSSEWithPolicies[I any](s ServerRuntime, method, path string, eventTypeMap map[string]any, handler SSEHandler[I], policies ...SSERoutePolicy[I]) {
	lo.Must0(RouteSSEWithPolicies(s, method, path, eventTypeMap, handler, policies...))
}

// MustGroupRouteSSEWithPolicies registers a grouped SSE route with policies and panics on failure.
func MustGroupRouteSSEWithPolicies[I any](g *Group, method, path string, eventTypeMap map[string]any, handler SSEHandler[I], policies ...SSERoutePolicy[I]) {
	lo.Must0(GroupRouteSSEWithPolicies(g, method, path, eventTypeMap, handler, policies...))
}
