package httpx

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/samber/lo"
)

// withInputValidation applies validator checks and standard error conversion.
func withInputValidation[I, O any](s *Server, handler TypedHandler[I, O]) TypedHandler[I, O] {
	if handler == nil || s == nil {
		return handler
	}

	validateInput := compileInputValidator[I](s.validator)

	return func(ctx context.Context, input *I) (out *O, err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				out = nil
				err = recoverTypedHandlerPanic(s, recovered)
			}
		}()

		if err := validateTypedInput(validateInput, input); err != nil {
			return nil, err
		}
		return invokeTypedHandler(ctx, handler, input)
	}
}

func recoverTypedHandlerPanic(s *Server, recovered any) error {
	if s == nil || !s.panicRecover {
		panic(recovered)
	}
	return recoverTypedHandlerError(recovered)
}

func invokeTypedHandler[I, O any](ctx context.Context, handler TypedHandler[I, O], input *I) (*O, error) {
	out, err := handler(ctx, input)
	if err != nil {
		return nil, convertTypedHandlerError(err)
	}
	return out, nil
}

func recoverTypedHandlerError(recovered any) error {
	return huma.Error500InternalServerError(fmt.Sprintf("panic in handler: %v", recovered))
}

func validateTypedInput[I any](validateInput func(*I) error, input *I) error {
	if validateInput == nil {
		return nil
	}
	if err := validateInput(input); err != nil {
		message := validationErrorMessage(err)
		return huma.Error400BadRequest(message, err)
	}
	return nil
}

func convertTypedHandlerError(err error) error {
	if httpxErr, ok := errors.AsType[*Error](err); ok {
		return lo.Ternary(
			httpxErr.Err != nil,
			huma.NewError(httpxErr.Code, httpxErr.Message, httpxErr.Err),
			huma.NewError(httpxErr.Code, httpxErr.Message),
		)
	}
	if statusErr, ok := errors.AsType[huma.StatusError](err); ok {
		return statusErr
	}
	return huma.Error500InternalServerError(err.Error(), err)
}
