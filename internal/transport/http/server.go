// internal/transport/http/server.go
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// Server represents HTTP server.
type Server struct {
	httpServer *http.Server
	port       int
	logger     logger.Logger
}

// NewServer creates new HTTP server.
func NewServer(handler http.Handler, port int, logger logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      handler,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		port:   port,
		logger: logger,
	}
}

// Start starts listening on configured port.
func (s *Server) Start() error {
	// TODO: Log server start on port
	return s.httpServer.ListenAndServe()
}

// Stop gracefully shuts down server.
func (s *Server) Stop(ctx context.Context) error {
	// TODO: Log server shutdown
	return s.httpServer.Shutdown(ctx)
}
