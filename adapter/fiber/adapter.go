package fiber

import (
	"github.com/arcgolabs/httpx/adapter"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
	"github.com/gofiber/fiber/v2"
)

// Adapter implements the fiber runtime bridge for httpx.
type Adapter struct {
	app  *fiber.App
	huma huma.API
}

// New constructs a fiber adapter backed by a fiber app and Huma API.
func New(app *fiber.App, opts ...adapter.HumaOptions) *Adapter {
	resolvedApp := orDefaultApp(app)
	humaOpts := adapter.MergeHumaOptions(opts...)
	cfg := huma.DefaultConfig(humaOpts.Title, humaOpts.Version)
	adapter.ApplyHumaConfig(&cfg, humaOpts)
	api := humafiber.New(resolvedApp, cfg)

	return &Adapter{
		app:  resolvedApp,
		huma: api,
	}
}

func orDefaultApp(app *fiber.App) *fiber.App {
	if app != nil {
		return app
	}
	return fiber.New()
}

// Name returns the adapter name.
func (a *Adapter) Name() string {
	return "fiber"
}

// Router exposes the underlying fiber app.
func (a *Adapter) Router() *fiber.App {
	return a.app
}

// HumaAPI exposes the underlying Huma API.
func (a *Adapter) HumaAPI() huma.API {
	return a.huma
}
