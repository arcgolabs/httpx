package fiber_test

import (
	"fmt"
	"net/http"

	fiberadapter "github.com/arcgolabs/httpx/adapter/fiber"
	fiberframework "github.com/gofiber/fiber/v3"
)

func testRequest(adapter *fiberadapter.Adapter, request *http.Request) (*http.Response, error) {
	response, err := adapter.Router().Test(request, fiberframework.TestConfig{Timeout: 0})
	if err != nil {
		return nil, fmt.Errorf("test fiber request: %w", err)
	}
	return response, nil
}
