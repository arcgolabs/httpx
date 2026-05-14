package httpx

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
)

var (
	// ErrInvalidETag reports a malformed HTTP ETag header value.
	ErrInvalidETag = errors.New("invalid etag")
	// ErrInvalidByteRange reports a malformed HTTP Range header value.
	ErrInvalidByteRange = errors.New("invalid byte range")
	// ErrInvalidContentMD5 reports a malformed Content-MD5 header value.
	ErrInvalidContentMD5 = errors.New("invalid content-md5")
)

func stringHeaderSchema(example string) *huma.Schema {
	return &huma.Schema{
		Type:     huma.TypeString,
		Examples: []any{example},
	}
}
