// internal/transport/http/server.go
package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// Server представляет HTTP-сервер.
type Server struct {
	httpServer *http.Server
	port       int
	logger     *logger.Logger
}

// NewServer создаёт новый HTTP-сервер.
func NewServer(handler http.Handler, port int, logger *logger.Logger) *Server {
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

// Start начинает слушание на настроенном порту.
func (s *Server) Start() error {
	s.logger.Info("starting HTTP server", map[string]interface{}{"port": s.port})
	return s.httpServer.ListenAndServe()
}

// Stop корректно завершает работу сервера.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping HTTP server", nil)
	return s.httpServer.Shutdown(ctx)
}
