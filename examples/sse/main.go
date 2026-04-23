// Package main demonstrates httpx server-sent events with typed messages.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/arcgolabs/httpx/examples/shared"
	"github.com/arcgolabs/httpx"
	"github.com/arcgolabs/httpx/adapter"
	"github.com/arcgolabs/httpx/adapter/std"
	"github.com/arcgolabs/pkg/randomport"
	"github.com/danielgtaylor/huma/v2"
)

type streamInput struct {
	Count int `query:"count"`
}

type tickEvent struct {
	Index int    `json:"index"`
	At    string `json:"at"`
}

type doneEvent struct {
	Message string `json:"message"`
}

func main() {
	logger, closeLogger, err := shared.NewLogger()
	if err != nil {
		panic(err)
	}

	a := std.New(nil, adapter.HumaOptions{
		Title:       "httpx SSE example",
		Version:     "1.0.0",
		Description: "SSE streaming with typed events",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	s := httpx.New(httpx.WithAdapter(a))

	httpx.MustRouteSSEWithPolicies(s, httpx.MethodGet, "/events", map[string]any{
		"tick": tickEvent{},
		"done": doneEvent{},
	}, streamEvents, httpx.SSEPolicyOperation[streamInput](huma.OperationTags("streaming")))

	port := randomport.MustFind()
	addr := fmt.Sprintf(":%d", port)
	logger.Info("example server starting",
		slog.String("example", "sse"),
		slog.String("address", addr),
		slog.String("openapi", fmt.Sprintf("http://localhost%s/openapi.json", addr)),
		slog.String("docs", fmt.Sprintf("http://localhost%s/docs", addr)),
		slog.String("stream_url", fmt.Sprintf("http://localhost%s/events?count=5", addr)),
		slog.String("curl", fmt.Sprintf("curl -N http://localhost%s/events?count=5", addr)),
	)

	if err := s.ListenPort(port); err != nil {
		logger.Error("server exited with error", slog.String("error", err.Error()))
		closeLogger()
		os.Exit(1)
	}
	closeLogger()
}

func streamEvents(ctx context.Context, input *streamInput, send httpx.SSESender) {
	count := input.Count
	if count <= 0 {
		count = 3
	}

	for i := 1; i <= count; i++ {
		if stopped := sendTick(ctx, send, i); stopped {
			return
		}
	}

	if err := send(httpx.SSEMessage{
		ID:   count + 1,
		Data: doneEvent{Message: "stream completed"},
	}); err != nil {
		return
	}
}

func sendTick(ctx context.Context, send httpx.SSESender, index int) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}

	if err := send(httpx.SSEMessage{
		ID:   index,
		Data: tickEvent{Index: index, At: time.Now().UTC().Format(time.RFC3339Nano)},
	}); err != nil {
		return true
	}
	return false
}
