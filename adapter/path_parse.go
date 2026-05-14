package adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

var errMultipleCatchAllParams = errors.New("multiple catch-all path parameters are not supported")

func catchAllNameFromOperation(op *huma.Operation) string {
	if op == nil {
		return ""
	}
	return catchAllName(op.Metadata)
}

func normalizeCatchAllPath(path string) (string, string, error) {
	if !strings.Contains(path, "...") {
		return path, "", nil
	}

	segments, leading := splitTemplatePath(path)
	catchAllName := ""
	for i, segment := range segments {
		name, matched, err := normalizeCatchAllSegment(segment, i == len(segments)-1)
		if err != nil {
			return "", "", err
		}
		if !matched {
			continue
		}
		if catchAllName != "" {
			return "", "", errMultipleCatchAllParams
		}
		catchAllName = name
		segments[i] = "{" + name + "}"
	}
	if catchAllName == "" {
		return "", "", fmt.Errorf("invalid catch-all path syntax in %q", path)
	}
	return joinTemplatePath(segments, leading), catchAllName, nil
}

func normalizeCatchAllSegment(segment string, final bool) (string, bool, error) {
	name, ok := parseCatchAllSegment(segment)
	if !ok {
		if strings.Contains(segment, "...") {
			return "", false, fmt.Errorf("invalid catch-all path segment %q", segment)
		}
		return "", false, nil
	}
	if !final {
		return "", false, fmt.Errorf("catch-all path parameter %q must be the last segment", segment)
	}
	return name, true, nil
}

func replaceCatchAllSegment(path, name, replacement string) (string, int, bool) {
	segments, leading := splitTemplatePath(path)
	target := "{" + name + "}"
	for i, segment := range segments {
		if segment != target {
			continue
		}
		segments[i] = replacement
		return joinTemplatePath(segments, leading), i, true
	}
	return path, -1, false
}

func splitTemplatePath(path string) ([]string, bool) {
	if path == "" {
		return nil, false
	}
	leading := strings.HasPrefix(path, "/")
	trimmed := strings.TrimPrefix(path, "/")
	if trimmed == "" {
		return nil, leading
	}
	return strings.Split(trimmed, "/"), leading
}

func joinTemplatePath(segments []string, leading bool) string {
	joined := strings.Join(segments, "/")
	if leading {
		return "/" + joined
	}
	return joined
}

func parseCatchAllSegment(segment string) (string, bool) {
	if !strings.HasPrefix(segment, "{") || !strings.HasSuffix(segment, "...}") {
		return "", false
	}
	name := strings.TrimSuffix(strings.TrimPrefix(segment, "{"), "...}")
	return name, name != ""
}
