// Package app provides a framework for building applications with a consistent lifecycle.
package app

import "context"

// Server defines the interface for a runnable service that can be managed by an App.
type Server interface {
	// Start starts the server. It is expected to be a non-blocking call,
	// often starting a new goroutine for the server's main loop.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the server, with a deadline provided by the context.
	Stop(ctx context.Context) error
}
