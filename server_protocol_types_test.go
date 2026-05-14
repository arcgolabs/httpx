package httpx_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_StrongTypedProtocolHeaders(t *testing.T) {
	server := newServer()

	type in struct {
		Range      ByteRange  `header:"Range"`
		IfMatch    ETagSet    `header:"If-Match"`
		ContentMD5 ContentMD5 `header:"Content-MD5"`
	}
	type out struct {
		Body struct {
			Range    string `json:"range"`
			Start    int64  `json:"start"`
			End      int64  `json:"end"`
			Matched  bool   `json:"matched"`
			Checksum string `json:"checksum"`
		}
	}

	err := Get(server, "/headers", func(_ context.Context, input *in) (*out, error) {
		start, end, ok := input.Range.Bounds(20)
		require.True(t, ok)

		result := &out{}
		result.Body.Range = input.Range.String()
		result.Body.Start = start
		result.Body.End = end
		result.Body.Matched = input.IfMatch.Contains(ETag{Value: "abc"})
		result.Body.Checksum = input.ContentMD5.String()
		return result, nil
	})
	require.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/headers", nil)
	req.Header.Set("Range", "bytes=5-9")
	req.Header.Set("If-Match", `"abc", W/"def"`)
	req.Header.Set("Content-MD5", "1B2M2Y8AsgTpgAmY7PhCfg==")
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"range":"bytes=5-9"`)
	assert.Contains(t, rec.Body.String(), `"start":5`)
	assert.Contains(t, rec.Body.String(), `"end":9`)
	assert.Contains(t, rec.Body.String(), `"matched":true`)
}

func TestServer_StrongTypedProtocolHeaderValidation(t *testing.T) {
	server := newServer()

	type in struct {
		Range ByteRange `header:"Range"`
	}

	err := Get(server, "/headers/invalid", func(_ context.Context, _ *in) (*pingOutput, error) {
		out := &pingOutput{}
		out.Body.Message = "ok"
		return out, nil
	})
	require.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/headers/invalid", nil)
	req.Header.Set("Range", "items=5-9")
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "Range")
}

func TestServer_RequestStreamOnTypedRoute(t *testing.T) {
	server := newServer()

	type in struct {
		ContentLength int64 `header:"Content-Length"`
		Payload       RequestStream
	}
	type out struct {
		Body struct {
			Size int64  `json:"size"`
			Text string `json:"text"`
		}
	}

	err := Put(server, "/upload", func(_ context.Context, input *in) (*out, error) {
		data, err := io.ReadAll(input.Payload.Reader())
		require.NoError(t, err)

		result := &out{}
		result.Body.Size = int64(len(data))
		result.Body.Text = string(data)
		assert.Equal(t, input.ContentLength, result.Body.Size)
		return result, nil
	}, OperationBinaryRequest("application/octet-stream"))
	require.NoError(t, err)

	req := newTestRequest(http.MethodPut, "/upload", bytes.NewReader([]byte("stream-body")))
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", "11")
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"size":11`)
	assert.Contains(t, rec.Body.String(), `"text":"stream-body"`)

	pathItem := server.OpenAPI().Paths["/upload"]
	if assert.NotNil(t, pathItem) && assert.NotNil(t, pathItem.Put) && assert.NotNil(t, pathItem.Put.RequestBody) {
		assert.Contains(t, pathItem.Put.RequestBody.Content, "application/octet-stream")
	}
}

func TestServer_ResponseStreamOnTypedRoute(t *testing.T) {
	server := newServer()

	type in struct {
		Name string `path:"name"`
	}
	type out struct {
		ContentType string `header:"Content-Type"`
		Body        ResponseStream
	}

	err := Get(server, "/download/{name}", func(_ context.Context, input *in) (*out, error) {
		return &out{
			ContentType: "text/plain",
			Body:        StreamReader(strings.NewReader("download:" + input.Name)),
		}, nil
	}, OperationBinaryResponse("text/plain"))
	require.NoError(t, err)

	req := newTestRequest(http.MethodGet, "/download/report", nil)
	rec := serveRequest(t, server, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/plain")
	assert.Equal(t, "download:report", rec.Body.String())
}
