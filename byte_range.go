package httpx

import (
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

// ByteRange represents a single bytes Range header value.
type ByteRange struct {
	Start        int64
	End          int64
	SuffixLength int64
	HasStart     bool
	HasEnd       bool
	Suffix       bool
}

// ParseByteRange parses a single bytes Range header value.
func ParseByteRange(value string) (ByteRange, error) {
	spec, err := parseByteRangeSpec(value)
	if err != nil {
		return ByteRange{}, err
	}

	startText, endText, ok := strings.Cut(spec, "-")
	if !ok || startText == "" && endText == "" {
		return ByteRange{}, ErrInvalidByteRange
	}

	switch {
	case startText == "":
		return parseSuffixByteRange(endText)
	case endText == "":
		return parseOpenByteRange(startText)
	default:
		return parseClosedByteRange(startText, endText)
	}
}

// UnmarshalText parses a byte range from a header value.
func (r *ByteRange) UnmarshalText(text []byte) error {
	parsed, err := ParseByteRange(string(text))
	if err != nil {
		return err
	}
	*r = parsed
	return nil
}

// String returns the canonical header representation.
func (r ByteRange) String() string {
	switch {
	case r.Suffix:
		return "bytes=-" + strconv.FormatInt(r.SuffixLength, 10)
	case r.HasStart && r.HasEnd:
		return "bytes=" + strconv.FormatInt(r.Start, 10) + "-" + strconv.FormatInt(r.End, 10)
	case r.HasStart:
		return "bytes=" + strconv.FormatInt(r.Start, 10) + "-"
	default:
		return ""
	}
}

// Bounds resolves the byte range against a known object size.
func (r ByteRange) Bounds(size int64) (start, end int64, ok bool) {
	if size <= 0 {
		return 0, 0, false
	}
	if r.Suffix {
		return suffixBounds(size, r.SuffixLength)
	}
	if !r.HasStart || r.Start >= size {
		return 0, 0, false
	}
	if !r.HasEnd || r.End >= size {
		return r.Start, size - 1, true
	}
	return r.Start, r.End, true
}

// IsZero reports whether r was not set.
func (r ByteRange) IsZero() bool {
	return !r.HasStart && !r.HasEnd && !r.Suffix
}

// Schema documents ByteRange as a string parameter.
func (r ByteRange) Schema(_ huma.Registry) *huma.Schema {
	return stringHeaderSchema("bytes=0-1023")
}

func parseByteRangeSpec(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	unit, spec, ok := strings.Cut(trimmed, "=")
	if !ok || !strings.EqualFold(strings.TrimSpace(unit), "bytes") {
		return "", ErrInvalidByteRange
	}
	spec = strings.TrimSpace(spec)
	if strings.Contains(spec, ",") {
		return "", ErrInvalidByteRange
	}
	return spec, nil
}

func parseSuffixByteRange(endText string) (ByteRange, error) {
	suffix, err := parseNonNegativeInt64(endText)
	if err != nil || suffix <= 0 {
		return ByteRange{}, ErrInvalidByteRange
	}
	return ByteRange{SuffixLength: suffix, Suffix: true}, nil
}

func parseOpenByteRange(startText string) (ByteRange, error) {
	start, err := parseNonNegativeInt64(startText)
	if err != nil {
		return ByteRange{}, ErrInvalidByteRange
	}
	return ByteRange{Start: start, HasStart: true}, nil
}

func parseClosedByteRange(startText, endText string) (ByteRange, error) {
	start, err := parseNonNegativeInt64(startText)
	if err != nil {
		return ByteRange{}, ErrInvalidByteRange
	}
	end, err := parseNonNegativeInt64(endText)
	if err != nil || end < start {
		return ByteRange{}, ErrInvalidByteRange
	}
	return ByteRange{Start: start, End: end, HasStart: true, HasEnd: true}, nil
}

func suffixBounds(size, suffixLength int64) (start, end int64, ok bool) {
	if suffixLength >= size {
		return 0, size - 1, true
	}
	return size - suffixLength, size - 1, true
}

func parseNonNegativeInt64(value string) (int64, error) {
	parsed, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64)
	if err != nil || parsed < 0 {
		return 0, ErrInvalidByteRange
	}
	return parsed, nil
}
