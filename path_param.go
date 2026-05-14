package httpx

import (
	"reflect"

	"github.com/danielgtaylor/huma/v2"
)

type rawPathParamContext interface {
	RawParam(name string) (string, bool)
}

type humaContextUnwrapper interface {
	Unwrap() huma.Context
}

// RawPathParam returns the escaped value for an httpx catch-all path parameter.
func RawPathParam(ctx huma.Context, name string) (string, bool) {
	for ctx != nil {
		if raw, ok := ctx.(rawPathParamContext); ok {
			return raw.RawParam(name)
		}
		unwrapper, ok := ctx.(humaContextUnwrapper)
		if !ok {
			return "", false
		}
		ctx = unwrapper.Unwrap()
	}
	return "", false
}

// PathTail binds a catch-all path parameter as decoded and escaped text.
type PathTail struct {
	Value    string `json:"-"`
	RawValue string `json:"-"`
	IsSet    bool   `json:"-"`
}

// Receiver exposes the decoded value destination to Huma's parameter parser.
func (p *PathTail) Receiver() reflect.Value {
	return reflect.ValueOf(p).Elem().FieldByName("Value")
}

// OnParamSet records whether Huma saw the path parameter.
func (p *PathTail) OnParamSet(isSet bool, parsed any) {
	p.IsSet = isSet
	if value, ok := parsed.(string); ok {
		p.Value = value
	}
}

// Resolve stores the escaped catch-all value after parameter parsing.
func (p *PathTail) Resolve(ctx huma.Context) []error {
	if raw, ok := RawPathParam(ctx, ""); ok {
		p.RawValue = raw
	}
	return nil
}

// Schema describes PathTail as a string in OpenAPI.
func (p PathTail) Schema(_ huma.Registry) *huma.Schema {
	return &huma.Schema{Type: huma.TypeString}
}

// String returns the decoded path tail.
func (p PathTail) String() string {
	return p.Value
}

// Raw returns the escaped path tail.
func (p PathTail) Raw() string {
	return p.RawValue
}

// Present reports whether the route parameter was present in the request.
func (p PathTail) Present() bool {
	return p.IsSet
}
