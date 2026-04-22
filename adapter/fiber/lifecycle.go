package fiber

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/samber/oops"
)

// Listen starts the fiber server.
func (a *Adapter) Listen(addr string) error {
	if err := a.app.Listen(addr); err != nil {
		return wrapFiberListenError(addr, err)
	}
	return nil
}

// Shutdown stops the fiber server.
func (a *Adapter) Shutdown() error {
	if err := a.app.Shutdown(); err != nil && !isExpectedFiberClose(err) {
		return oops.In("httpx/adapter/fiber").
			With("op", "shutdown").
			Wrapf(err, "httpx/fiber: shutdown")
	}
	return nil
}

// ListenContext starts related services.
func (a *Adapter) ListenContext(ctx context.Context, addr string) error {
	if ctx == nil {
		return a.Listen(addr)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- a.Listen(addr)
	}()

	select {
	case err := <-errCh:
		if isExpectedFiberClose(err) {
			return nil
		}
		return wrapFiberListenError(addr, err)
	case <-ctx.Done():
		shutdownErr := a.Shutdown()
		listenErr := <-errCh
		if shutdownErr != nil {
			return oops.In("httpx/adapter/fiber").
				With("op", "shutdown", "addr", addr).
				Wrapf(shutdownErr, "httpx/fiber: shutdown on %q", addr)
		}
		if isExpectedFiberClose(listenErr) {
			return nil
		}
		return wrapFiberListenError(addr, listenErr)
	}
}

func wrapFiberListenError(addr string, err error) error {
	return oops.In("httpx/adapter/fiber").
		With("op", "listen", "addr", addr).
		Wrapf(err, "httpx/fiber: listen on %q", addr)
}

func isExpectedFiberClose(err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, http.ErrServerClosed) {
		return true
	}
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "server is not running") ||
		strings.Contains(lower, "use of closed network connection")
}
