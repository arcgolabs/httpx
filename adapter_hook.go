package httpx

import "github.com/arcgolabs/httpx/adapter"

func useHostCapability[T any](s ServerRuntime, use func(T)) bool {
	server := unwrapServer(s)
	if server == nil || server.adapter == nil || use == nil {
		return false
	}
	return withHostCapability(server.adapter, use)
}

func withHostCapability[T any](a adapter.Host, use func(T)) bool {
	if a == nil || use == nil {
		return false
	}

	capability, ok := any(a).(T)
	if !ok {
		return false
	}
	use(capability)
	return true
}
