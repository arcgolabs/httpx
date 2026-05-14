//go:build !no_gin

package gin

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/arcgolabs/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

type lifecycleState struct {
	mu     sync.Mutex
	server *http.Server
}

// Adapter implements the gin runtime bridge for httpx.
type Adapter struct {
	engine    *gin.Engine
	huma      huma.API
	lifecycle *lifecycleState
}

// New constructs a gin adapter backed by a gin engine and Huma API.
func New(engine *gin.Engine, opts ...adapter.HumaOptions) *Adapter {
	eng := orDefaultEngine(engine)
	humaOpts := adapter.MergeHumaOptions(opts...)
	cfg := huma.DefaultConfig(humaOpts.Title, humaOpts.Version)
	adapter.ApplyHumaConfig(&cfg, humaOpts)
	baseAPI := humagin.New(eng, adapter.RouterOnlyConfig(cfg))
	api := huma.NewAPI(cfg, adapter.NewPathAdapter(baseAPI.Adapter(), adapter.RouterPathGin))

	return &Adapter{
		engine:    eng,
		huma:      api,
		lifecycle: &lifecycleState{},
	}
}

func orDefaultEngine(engine *gin.Engine) *gin.Engine {
	if engine != nil {
		return engine
	}
	return gin.New()
}

// Name returns the adapter name.
func (a *Adapter) Name() string {
	return "gin"
}

// Router exposes the underlying gin engine.
func (a *Adapter) Router() *gin.Engine {
	return a.engine
}

// Listen starts related services.
func (a *Adapter) Listen(addr string) error {
	server := a.httpServer(addr)
	release := a.trackServer(server)
	defer release()

	if err := server.ListenAndServe(); err != nil {
		return wrapGinListenError(addr, err)
	}
	return nil
}

// Shutdown stops the active gin server.
func (a *Adapter) Shutdown() error {
	return a.shutdownContext(context.Background())
}

func (a *Adapter) shutdownContext(ctx context.Context) error {
	server := a.activeServer()
	if server == nil {
		return nil
	}
	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return oops.In("httpx/adapter/gin").
			With("op", "shutdown", "addr", server.Addr).
			Wrapf(err, "httpx/gin: shutdown")
	}
	return nil
}

// ListenContext starts related services.
func (a *Adapter) ListenContext(ctx context.Context, addr string) error {
	if ctx == nil {
		return a.Listen(addr)
	}

	server := a.httpServer(addr)
	release := a.trackServer(server)
	defer release()

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return wrapGinListenError(addr, err)
	case <-ctx.Done():
		if err := a.shutdownContext(ctx); err != nil {
			return oops.In("httpx/adapter/gin").
				With("op", "shutdown", "addr", addr).
				Wrapf(err, "httpx/gin: shutdown on %q", addr)
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return wrapGinListenError(addr, err)
	}
}

// HumaAPI exposes the underlying Huma API.
func (a *Adapter) HumaAPI() huma.API {
	return a.huma
}

func (a *Adapter) httpServer(addr string) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           a.engine.Handler(),
		ReadHeaderTimeout: 30 * time.Second,
	}
}

func (a *Adapter) trackServer(server *http.Server) func() {
	if a == nil || a.lifecycle == nil {
		return func() {}
	}

	a.lifecycle.mu.Lock()
	a.lifecycle.server = server
	a.lifecycle.mu.Unlock()

	return func() {
		a.lifecycle.mu.Lock()
		if a.lifecycle.server == server {
			a.lifecycle.server = nil
		}
		a.lifecycle.mu.Unlock()
	}
}

func (a *Adapter) activeServer() *http.Server {
	if a == nil || a.lifecycle == nil {
		return nil
	}

	a.lifecycle.mu.Lock()
	defer a.lifecycle.mu.Unlock()
	return a.lifecycle.server
}

func wrapGinListenError(addr string, err error) error {
	return oops.In("httpx/adapter/gin").
		With("op", "listen", "addr", addr).
		Wrapf(err, "httpx/gin: listen on %q", addr)
}
