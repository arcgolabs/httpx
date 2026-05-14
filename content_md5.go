package httpx

import (
	"encoding/base64"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// ContentMD5 represents a validated Content-MD5 header value.
type ContentMD5 struct {
	sum [16]byte
	raw string
}

// ParseContentMD5 parses and validates one Content-MD5 header value.
func ParseContentMD5(value string) (ContentMD5, error) {
	trimmed := strings.TrimSpace(value)
	decoded, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil || len(decoded) != 16 {
		return ContentMD5{}, ErrInvalidContentMD5
	}

	result := ContentMD5{raw: trimmed}
	copy(result.sum[:], decoded)
	return result, nil
}

// UnmarshalText parses Content-MD5 from a header value.
func (m *ContentMD5) UnmarshalText(text []byte) error {
	parsed, err := ParseContentMD5(string(text))
	if err != nil {
		return err
	}
	*m = parsed
	return nil
}

// Bytes returns a copy of the decoded MD5 digest.
func (m ContentMD5) Bytes() []byte {
	out := make([]byte, len(m.sum))
	copy(out, m.sum[:])
	return out
}

// String returns the canonical header representation.
func (m ContentMD5) String() string {
	if m.raw != "" {
		return m.raw
	}
	if m.IsZero() {
		return ""
	}
	return base64.StdEncoding.EncodeToString(m.sum[:])
}

// IsZero reports whether m was not set.
func (m ContentMD5) IsZero() bool {
	return m.raw == "" && m.sum == [16]byte{}
}

// Schema documents ContentMD5 as a string parameter.
func (m ContentMD5) Schema(_ huma.Registry) *huma.Schema {
	return stringHeaderSchema("1B2M2Y8AsgTpgAmY7PhCfg==")
}
