package httpx

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
)

// compileInputValidator compiles one typed validation function so request path does not
// repeat reflection shape checks for every invocation.
func compileInputValidator[I any](v *validator.Validate) func(*I) error {
	if v == nil {
		return nil
	}

	_, hasNestedPointer, ok := indirectStructType[I]()
	if !ok {
		return nil
	}

	if !hasNestedPointer {
		return directInputValidator[I](v)
	}

	return nestedPointerInputValidator[I](v)
}

func directInputValidator[I any](v *validator.Validate) func(*I) error {
	return func(input *I) error {
		if input == nil {
			return nil
		}
		return v.Struct(input)
	}
}

func nestedPointerInputValidator[I any](v *validator.Validate) func(*I) error {
	return func(input *I) error {
		if _, ok := indirectStructValue(input); !ok {
			return nil
		}
		return v.Struct(input)
	}
}

// validationErrorMessage converts validator errors into a concise HTTP-facing message.
func validationErrorMessage(err error) string {
	var validationErrs validator.ValidationErrors
	if !errors.As(err, &validationErrs) {
		return "request validation failed"
	}

	issues := lo.Map(validationErrs, func(validationErr validator.FieldError, _ int) string {
		field := validationErr.Field()
		if field == "" {
			field = validationErr.StructField()
		}
		if field == "" {
			field = "input"
		}

		return field + " failed '" + validationErr.Tag() + "'"
	})

	if len(issues) == 0 {
		return "request validation failed"
	}

	return "request validation failed: " + strings.Join(issues, "; ")
}
