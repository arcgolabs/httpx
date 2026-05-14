package httpx

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/arcgolabs/collectionx/set"
	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
)

// OperationBinaryResponse documents a binary payload for HTTP 200.
func OperationBinaryResponse(contentTypes ...string) OperationOption {
	normalized := normalizeContentTypes(contentTypes, "application/octet-stream")

	return func(op *huma.Operation) {
		if op == nil {
			return
		}
		response := ensureHTTPResponse(op, http.StatusOK)
		appendBinaryContentTypes(response, normalized)
	}
}

// OperationBinaryRequest documents a binary request payload.
func OperationBinaryRequest(contentTypes ...string) OperationOption {
	normalized := normalizeContentTypes(contentTypes, "application/octet-stream")

	return func(op *huma.Operation) {
		if op == nil {
			return
		}
		ensureBinaryRequestBody(op, normalized)
	}
}

// OperationHTMLResponse documents HTML payload for HTTP 200.
func OperationHTMLResponse() OperationOption {
	return func(op *huma.Operation) {
		if op == nil {
			return
		}
		response := ensureHTTPResponse(op, http.StatusOK)
		if _, exists := response.Content["text/html"]; exists {
			return
		}
		response.Content["text/html"] = &huma.MediaType{
			Schema: &huma.Schema{
				Type: huma.TypeString,
			},
		}
	}
}

func ensureBinaryRequestBody(op *huma.Operation, contentTypes []string) {
	if op.RequestBody == nil {
		op.RequestBody = &huma.RequestBody{
			Required: true,
		}
	}
	if op.RequestBody.Content == nil {
		op.RequestBody.Content = map[string]*huma.MediaType{}
	}
	lo.ForEach(contentTypes, func(contentType string, _ int) {
		if _, exists := op.RequestBody.Content[contentType]; exists {
			return
		}
		op.RequestBody.Content[contentType] = &huma.MediaType{
			Schema: &huma.Schema{
				Type:   huma.TypeString,
				Format: "binary",
			},
		}
	})
}

// PolicyImageResponse applies runtime default Content-Type and OpenAPI binary response.
func PolicyImageResponse[I, O any](contentTypes ...string) RoutePolicy[I, O] {
	normalized := normalizeContentTypes(contentTypes, "image/png")
	defaultType := normalized[0]
	headerSetter := compileHeaderSetter[O]("Content-Type", defaultType)

	return RoutePolicy[I, O]{
		Name:      "image-response",
		Operation: OperationBinaryResponse(normalized...),
		Wrap: func(next TypedHandler[I, O]) TypedHandler[I, O] {
			if next == nil {
				return nil
			}
			return func(ctx context.Context, input *I) (*O, error) {
				out, err := next(ctx, input)
				if err != nil || out == nil {
					return out, err
				}
				if headerSetter != nil {
					headerSetter(out)
				}
				return out, nil
			}
		},
	}
}

// PolicyHTMLResponse applies runtime default Content-Type and OpenAPI HTML response.
func PolicyHTMLResponse[I, O any]() RoutePolicy[I, O] {
	headerSetter := compileHeaderSetter[O]("Content-Type", "text/html")

	return RoutePolicy[I, O]{
		Name:      "html-response",
		Operation: OperationHTMLResponse(),
		Wrap: func(next TypedHandler[I, O]) TypedHandler[I, O] {
			if next == nil {
				return nil
			}
			return func(ctx context.Context, input *I) (*O, error) {
				out, err := next(ctx, input)
				if err != nil || out == nil {
					return out, err
				}
				if headerSetter != nil {
					headerSetter(out)
				}
				return out, nil
			}
		},
	}
}

func compileHeaderSetter[O any](headerName, headerValue string) func(*O) {
	if headerName == "" || headerValue == "" {
		return nil
	}

	outputType, _, ok := indirectStructType[O]()
	if !ok {
		return nil
	}

	fieldIndex, ok := headerFieldIndex(outputType, headerName)
	if !ok {
		return nil
	}

	return func(output *O) {
		setHeaderField(output, fieldIndex, headerValue)
	}
}

func ensureHTTPResponse(op *huma.Operation, status int) *huma.Response {
	code := strconv.Itoa(status)
	if op.Responses == nil {
		op.Responses = map[string]*huma.Response{}
	}
	if op.Responses[code] == nil {
		op.Responses[code] = &huma.Response{
			Description: http.StatusText(status),
		}
	}
	if op.Responses[code].Content == nil {
		op.Responses[code].Content = map[string]*huma.MediaType{}
	}
	return op.Responses[code]
}

func appendBinaryContentTypes(response *huma.Response, contentTypes []string) {
	lo.ForEach(contentTypes, func(contentType string, _ int) {
		if _, exists := response.Content[contentType]; exists {
			return
		}
		response.Content[contentType] = &huma.MediaType{
			Schema: &huma.Schema{
				Type:   huma.TypeString,
				Format: "binary",
			},
		}
	})
}

func headerFieldIndex(outputType reflect.Type, headerName string) (int, bool) {
	index, ok := lo.Find(lo.Range(outputType.NumField()), func(index int) bool {
		structField := outputType.Field(index)
		return strings.EqualFold(structField.Tag.Get("header"), headerName) && structField.Type.Kind() == reflect.String
	})
	if !ok {
		return 0, false
	}
	return index, true
}

func setHeaderField[O any](output *O, fieldIndex int, headerValue string) {
	value, ok := indirectStructValue(output)
	if !ok || fieldIndex >= value.NumField() {
		return
	}

	field := value.Field(fieldIndex)
	if field.Kind() == reflect.String && field.CanSet() && field.String() == "" {
		field.SetString(headerValue)
	}
}

func normalizeContentTypes(values []string, fallback string) []string {
	if len(values) == 0 {
		return []string{fallback}
	}

	ordered := set.NewOrderedSetWithCapacity[string](len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			ordered.Add(trimmed)
		}
	}
	normalized := ordered.Values()
	if len(normalized) == 0 {
		return []string{fallback}
	}
	return normalized
}
