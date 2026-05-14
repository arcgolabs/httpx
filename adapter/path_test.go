package adapter_test

import (
	"testing"

	"github.com/arcgolabs/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeOperationPath_CatchAll(t *testing.T) {
	op := &huma.Operation{
		Method: "GET",
		Path:   "/{bucket}/{key...}",
	}

	require.NoError(t, adapter.NormalizeOperationPath(op))

	assert.Equal(t, "/{bucket}/{key}", op.Path)
	compiled := adapter.CompileRouterOperationPath(op, adapter.RouterPathChi)
	require.NotNil(t, compiled.CatchAll)
	assert.Equal(t, "key", compiled.CatchAll.Name)
}

func TestCompileRouterOperationPath_CatchAll(t *testing.T) {
	tests := []struct {
		name       string
		kind       adapter.RouterPathKind
		routerPath string
		nativeName string
	}{
		{name: "chi", kind: adapter.RouterPathChi, routerPath: "/api/{bucket}/*", nativeName: "*"},
		{name: "gin", kind: adapter.RouterPathGin, routerPath: "/api/{bucket}/*key", nativeName: "key"},
		{name: "echo", kind: adapter.RouterPathEcho, routerPath: "/api/{bucket}/*", nativeName: "*"},
		{name: "fiber", kind: adapter.RouterPathFiber, routerPath: "/api/{bucket}/*", nativeName: "*"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := &huma.Operation{
				Method: "GET",
				Path:   "/api/{bucket}/{key...}",
			}
			require.NoError(t, adapter.NormalizeOperationPath(op))

			compiled := adapter.CompileRouterOperationPath(op, tt.kind)

			require.NotNil(t, compiled.CatchAll)
			assert.Equal(t, tt.routerPath, compiled.RouterPath)
			assert.Equal(t, tt.nativeName, compiled.CatchAll.NativeName)
			assert.Equal(t, 2, compiled.CatchAll.Index)
		})
	}
}

func TestNormalizeOperationPath_RejectsNonFinalCatchAll(t *testing.T) {
	op := &huma.Operation{
		Method: "GET",
		Path:   "/files/{key...}/{id}",
	}

	assert.Error(t, adapter.NormalizeOperationPath(op))
}
