package httpx

import (
	"context"
	"errors"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// PolicyTimeout applies a cooperative timeout to a route handler.
// The handler must respect context cancellation for the timeout to take effect.
func PolicyTimeout[I, O any](timeout time.Duration) RoutePolicy[I, O] {
	if timeout <= 0 {
		return RoutePolicy[I, O]{Name: "timeout"}
	}

	return RoutePolicy[I, O]{
		Name: "timeout",
		Wrap: func(next TypedHandler[I, O]) TypedHandler[I, O] {
			if next == nil {
				return nil
			}

			return func(ctx context.Context, input *I) (*O, error) {
				timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
				defer cancel()

				out, err := next(timeoutCtx, input)
				switch {
				case errors.Is(err, context.DeadlineExceeded),
					errors.Is(timeoutCtx.Err(), context.DeadlineExceeded),
					errors.Is(context.Cause(timeoutCtx), context.DeadlineExceeded):
					return nil, huma.Error504GatewayTimeout("request timeout", context.DeadlineExceeded)
				default:
					return out, err
				}
			}
		},
	}
}
