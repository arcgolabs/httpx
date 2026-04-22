// Package main demonstrates combining httpx HTTP routes with a websocket echo endpoint.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/DaiYuANg/arcgo/examples/httpx/shared"
	"github.com/DaiYuANg/arcgo/pkg/randomport"
	"github.com/arcgolabs/httpx"
	"github.com/arcgolabs/httpx/adapter"
	"github.com/arcgolabs/httpx/adapter/std"
	"github.com/arcgolabs/httpx/websocket"
)

type healthOutput struct {
	Body struct {
		Status string `json:"status"`
	} `json:"body"`
}

func main() {
	logger, closeLogger, err := shared.NewLogger()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	a := std.New(router, adapter.HumaOptions{
		Title:       "httpx websocket example",
		Version:     "1.0.0",
		Description: "Typed HTTP routes + websocket echo endpoint",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})
	s := httpx.New(httpx.WithAdapter(a))

	httpx.MustGet(s, "/health", func(_ context.Context, _ *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})

	router.HandleFunc("/ws/echo", websocket.HandlerFunc(echoWebSocket, websocket.WithCompression(true)))

	port := randomport.MustFind()
	addr := fmt.Sprintf(":%d", port)
	logger.Info("example server starting",
		slog.String("example", "websocket"),
		slog.String("address", addr),
		slog.String("openapi", fmt.Sprintf("http://localhost%s/openapi.json", addr)),
		slog.String("docs", fmt.Sprintf("http://localhost%s/docs", addr)),
		slog.String("health", fmt.Sprintf("http://localhost%s/health", addr)),
		slog.String("ws", fmt.Sprintf("ws://localhost%s/ws/echo", addr)),
	)

	if err := s.ListenPort(port); err != nil {
		logger.Error("server exited with error", slog.String("error", err.Error()))
		closeLogger()
		os.Exit(1)
	}
	closeLogger()
}

func echoWebSocket(ctx context.Context, conn websocket.Conn) error {
	for {
		msg, err := conn.Read(ctx)
		if err != nil {
			return fmt.Errorf("read websocket message: %w", err)
		}
		if err := conn.Write(msg); err != nil {
			return fmt.Errorf("write websocket message: %w", err)
		}
	}
}
