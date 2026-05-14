package httpx

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
)

// RequestStream injects the request body reader into a typed input without
// forcing the body into memory.
type RequestStream struct {
	ctx    huma.Context
	reader io.Reader
}

// Resolve captures the Huma request context after typed parameters are parsed.
func (s *RequestStream) Resolve(ctx huma.Context) []error {
	s.ctx = ctx
	if ctx == nil {
		s.reader = http.NoBody
		return nil
	}
	s.reader = ctx.BodyReader()
	if s.reader == nil {
		s.reader = http.NoBody
	}
	return nil
}

// Reader returns the request body reader.
func (s RequestStream) Reader() io.Reader {
	if s.reader == nil {
		return http.NoBody
	}
	return s.reader
}

// Context returns the request context.
func (s RequestStream) Context() context.Context {
	if s.ctx == nil {
		return context.Background()
	}
	return s.ctx.Context()
}

// URL returns the request URL captured by the underlying adapter.
func (s RequestStream) URL() url.URL {
	if s.ctx == nil {
		return url.URL{}
	}
	return s.ctx.URL()
}

// Header returns one request header value.
func (s RequestStream) Header(name string) string {
	if s.ctx == nil {
		return ""
	}
	return s.ctx.Header(name)
}

// ResponseStream is the typed output body callback used for streaming
// responses.
type ResponseStream = func(huma.Context)

// StreamWriter adapts a writer callback into a ResponseStream.
func StreamWriter(write func(io.Writer)) ResponseStream {
	return func(ctx huma.Context) {
		if ctx == nil || write == nil {
			return
		}
		write(ctx.BodyWriter())
	}
}

// StreamReader copies reader into the response body.
func StreamReader(reader io.Reader) ResponseStream {
	return StreamWriter(func(writer io.Writer) {
		if reader == nil || writer == nil {
			return
		}
		if _, err := io.Copy(writer, reader); err != nil {
			return
		}
	})
}

// StreamBytes writes bytes into the response body.
func StreamBytes(data []byte) ResponseStream {
	return StreamWriter(func(writer io.Writer) {
		if writer == nil {
			return
		}
		if _, err := writer.Write(data); err != nil {
			return
		}
	})
}
