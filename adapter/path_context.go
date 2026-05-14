package adapter

import (
	"net/url"
	"strings"
	"sync"

	"github.com/danielgtaylor/huma/v2"
)

// CatchAllRemainderPresent reports whether the request path includes the
// catch-all segment. `/bucket/` includes an empty segment, while `/bucket` does
// not.
func CatchAllRemainderPresent(ctx huma.Context, compiled CompiledRoutePath) bool {
	_, ok := catchAllRawValue(ctx, compiled)
	return ok
}

// WrapCatchAllContext returns a Huma context that exposes the logical catch-all
// parameter name via Param.
func WrapCatchAllContext(ctx huma.Context, op *huma.Operation, compiled CompiledRoutePath) huma.Context {
	if ctx == nil || compiled.CatchAll == nil {
		return ctx
	}
	return &catchAllContext{
		embeddedHumaContext: ctx,
		op:                  op,
		compiled:            compiled,
	}
}

type catchAllContext struct {
	embeddedHumaContext
	op                *huma.Operation
	compiled          CompiledRoutePath
	rawOnce           sync.Once
	rawValueCached    string
	rawValuePresent   bool
	paramOnce         sync.Once
	paramValueCached  string
	paramValuePresent bool
}

type embeddedHumaContext interface {
	huma.Context
}

func (c *catchAllContext) Operation() *huma.Operation {
	if c.op != nil {
		return c.op
	}
	return c.embeddedHumaContext.Operation()
}

func (c *catchAllContext) Param(name string) string {
	if c.compiled.CatchAll == nil || name != c.compiled.CatchAll.Name {
		return c.embeddedHumaContext.Param(name)
	}

	decoded, ok := c.decodedCatchAllValue()
	if !ok {
		return ""
	}
	return decoded
}

// RawParam returns the escaped catch-all parameter value. Passing an empty name
// returns the only catch-all parameter supported by httpx.
func (c *catchAllContext) RawParam(name string) (string, bool) {
	if c.compiled.CatchAll == nil {
		return "", false
	}
	if name != "" && name != c.compiled.CatchAll.Name {
		return "", false
	}
	return c.rawCatchAllValue()
}

func (c *catchAllContext) Unwrap() huma.Context {
	return c.embeddedHumaContext
}

func (c *catchAllContext) rawCatchAllValue() (string, bool) {
	c.rawOnce.Do(func() {
		c.rawValueCached, c.rawValuePresent = catchAllRawValue(c.embeddedHumaContext, c.compiled)
	})
	return c.rawValueCached, c.rawValuePresent
}

func (c *catchAllContext) decodedCatchAllValue() (string, bool) {
	c.paramOnce.Do(func() {
		raw, ok := c.rawCatchAllValue()
		if !ok {
			return
		}
		decoded, err := url.PathUnescape(raw)
		if err != nil {
			c.paramValueCached = raw
			c.paramValuePresent = true
			return
		}
		c.paramValueCached = decoded
		c.paramValuePresent = true
	})
	return c.paramValueCached, c.paramValuePresent
}

func catchAllRawValue(ctx huma.Context, compiled CompiledRoutePath) (string, bool) {
	if ctx == nil || compiled.CatchAll == nil || compiled.CatchAll.Index < 0 {
		return "", false
	}

	requestURL := ctx.URL()
	escapedPath := requestURL.EscapedPath()
	trimmed := strings.TrimPrefix(escapedPath, "/")
	if trimmed == "" {
		return "", false
	}

	segments := strings.Split(trimmed, "/")
	index := compiled.CatchAll.Index
	if len(segments) < index {
		return "", false
	}
	if len(segments) == index {
		return "", strings.HasSuffix(escapedPath, "/")
	}
	return strings.Join(segments[index:], "/"), true
}
