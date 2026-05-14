package httpx

import (
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// ETag represents one HTTP entity tag.
type ETag struct {
	Value    string
	Weak     bool
	Wildcard bool
}

// ParseETag parses one HTTP entity tag.
func ParseETag(value string) (ETag, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "*" {
		return ETag{Wildcard: true}, nil
	}

	weak := false
	if strings.HasPrefix(trimmed, "W/") {
		weak = true
		trimmed = strings.TrimSpace(strings.TrimPrefix(trimmed, "W/"))
	}
	if len(trimmed) < 2 || trimmed[0] != '"' || trimmed[len(trimmed)-1] != '"' {
		return ETag{}, ErrInvalidETag
	}
	inner, err := strconv.Unquote(trimmed)
	if err != nil || strings.ContainsAny(inner, "\r\n") {
		return ETag{}, ErrInvalidETag
	}
	return ETag{Value: inner, Weak: weak}, nil
}

// UnmarshalText parses an entity tag from a header value.
func (e *ETag) UnmarshalText(text []byte) error {
	parsed, err := ParseETag(string(text))
	if err != nil {
		return err
	}
	*e = parsed
	return nil
}

// String returns the canonical header representation.
func (e ETag) String() string {
	if e.Wildcard {
		return "*"
	}
	if e.IsZero() {
		return ""
	}
	prefix := ""
	if e.Weak {
		prefix = "W/"
	}
	return prefix + strconv.Quote(e.Value)
}

// Match reports whether e matches candidate.
func (e ETag) Match(candidate ETag) bool {
	return e.Wildcard || !e.IsZero() && e.Value == candidate.Value
}

// IsZero reports whether e was not set.
func (e ETag) IsZero() bool {
	return e.Value == "" && !e.Weak && !e.Wildcard
}

// Schema documents ETag as a string parameter.
func (e ETag) Schema(_ huma.Registry) *huma.Schema {
	return stringHeaderSchema(`"686897696a7c876b7e"`)
}

// ETagSet represents an HTTP entity tag list, such as If-Match.
type ETagSet []ETag

// ParseETagSet parses an HTTP entity tag list.
func ParseETagSet(value string) (ETagSet, error) {
	parts, err := splitHeaderList(value)
	if err != nil {
		return nil, err
	}
	if len(parts) == 0 {
		return nil, ErrInvalidETag
	}

	values := make(ETagSet, 0, len(parts))
	for _, part := range parts {
		etag, err := ParseETag(part)
		if err != nil {
			return nil, err
		}
		values = append(values, etag)
	}
	return values, nil
}

// UnmarshalText parses an entity tag list from a header value.
func (s *ETagSet) UnmarshalText(text []byte) error {
	parsed, err := ParseETagSet(string(text))
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// String returns the canonical header representation.
func (s ETagSet) String() string {
	parts := make([]string, 0, len(s))
	for _, etag := range s {
		parts = append(parts, etag.String())
	}
	return strings.Join(parts, ", ")
}

// Contains reports whether candidate is present in the set.
func (s ETagSet) Contains(candidate ETag) bool {
	for _, etag := range s {
		if etag.Match(candidate) {
			return true
		}
	}
	return false
}

// Schema documents ETagSet as a string parameter.
func (s ETagSet) Schema(_ huma.Registry) *huma.Schema {
	return stringHeaderSchema(`"686897696a7c876b7e", "726897696a7c876b7e"`)
}

type headerListState struct {
	parts   []string
	start   int
	quoted  bool
	escaped bool
}

func splitHeaderList(value string) ([]string, error) {
	state := &headerListState{}
	for index, char := range value {
		if err := state.consume(value, index, char); err != nil {
			return nil, err
		}
	}
	return state.finish(value)
}

func (s *headerListState) consume(value string, index int, char rune) error {
	if s.escaped {
		s.escaped = false
		return nil
	}

	switch char {
	case '\\':
		s.escaped = true
	case '"':
		s.quoted = !s.quoted
	case ',':
		if !s.quoted {
			return s.appendPart(value, index)
		}
	}
	return nil
}

func (s *headerListState) appendPart(value string, end int) error {
	part := strings.TrimSpace(value[s.start:end])
	if part == "" {
		return ErrInvalidETag
	}
	s.parts = append(s.parts, part)
	s.start = end + 1
	return nil
}

func (s *headerListState) finish(value string) ([]string, error) {
	if s.quoted || s.escaped {
		return nil, ErrInvalidETag
	}
	part := strings.TrimSpace(value[s.start:])
	if part == "" {
		return nil, ErrInvalidETag
	}
	return append(s.parts, part), nil
}
