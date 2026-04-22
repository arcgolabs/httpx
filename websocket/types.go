package websocket

import "errors"

var (
	// ErrClosed reports that the WebSocket connection is already closed.
	ErrClosed = errors.New("httpx/websocket: connection closed")
	// ErrUpgradeFailed reports that the WebSocket upgrade failed.
	ErrUpgradeFailed = errors.New("httpx/websocket: upgrade failed")
)

// MessageType identifies the type of a WebSocket message.
type MessageType uint8

const (
	// MessageText represents a text frame.
	MessageText MessageType = iota + 1
	// MessageBinary represents a binary frame.
	MessageBinary
	// MessagePing represents a ping control frame.
	MessagePing
	// MessagePong represents a pong control frame.
	MessagePong
	// MessageClose represents a close control frame.
	MessageClose
)

// Message wraps a typed WebSocket payload.
type Message struct {
	Type MessageType
	Data []byte
}

// Handler handles a WebSocket connection after upgrade.
type Handler func(Context, Conn) error

// Conn exposes typed read and write operations for a WebSocket connection.
type Conn interface {
	Read(Context) (Message, error)
	Write(Message) error
	Close(code uint16, reason []byte) error
}
