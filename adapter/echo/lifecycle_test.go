package echo_test

import (
	"context"
	"net"
	"testing"
	"time"

	echoadapter "github.com/arcgolabs/httpx/adapter/echo"
	"github.com/stretchr/testify/require"
)

const lifecycleTestTimeout = 3 * time.Second

func TestAdapterListen_ShutdownStopsServer(t *testing.T) {
	a := echoadapter.New(nil)
	addr := availableTCPAddress(t)
	errCh := make(chan error, 1)

	go func() {
		errCh <- a.Listen(addr)
	}()

	waitForTCPServer(t, addr)
	require.NoError(t, a.Shutdown())
	require.NoError(t, waitForListenResult(t, errCh))
}

func TestAdapterListenContext_CancellationStopsServer(t *testing.T) {
	a := echoadapter.New(nil)
	addr := availableTCPAddress(t)
	ctx, cancel := context.WithCancel(t.Context())
	errCh := make(chan error, 1)

	go func() {
		errCh <- a.ListenContext(ctx, addr)
	}()

	waitForTCPServer(t, addr)
	cancel()
	require.NoError(t, waitForListenResult(t, errCh))
}

func TestAdapterListen_OccupiedAddressReturnsError(t *testing.T) {
	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, listener.Close()) })

	a := echoadapter.New(nil)
	require.Error(t, a.Listen(listener.Addr().String()))
}

func availableTCPAddress(t *testing.T) string {
	t.Helper()

	listener, err := (&net.ListenConfig{}).Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	addr := listener.Addr().String()
	require.NoError(t, listener.Close())
	return addr
}

func waitForTCPServer(t *testing.T, addr string) {
	t.Helper()

	deadline := time.Now().Add(lifecycleTestTimeout)
	dialer := net.Dialer{Timeout: 50 * time.Millisecond}
	for time.Now().Before(deadline) {
		conn, err := dialer.DialContext(t.Context(), "tcp", addr)
		if err == nil {
			require.NoError(t, conn.Close())
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("TCP server %s did not start within %s", addr, lifecycleTestTimeout)
}

func waitForListenResult(t *testing.T, errCh <-chan error) error {
	t.Helper()

	select {
	case err := <-errCh:
		return err
	case <-time.After(lifecycleTestTimeout):
		t.Fatalf("listen did not return within %s", lifecycleTestTimeout)
		return nil
	}
}
