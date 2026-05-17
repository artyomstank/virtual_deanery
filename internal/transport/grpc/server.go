// internal/transport/grpc/server.go
package grpc

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/transport/grpc/handler"
	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"google.golang.org/grpc"
)

// Server represents gRPC server.
type Server struct {
	grpcServer *grpc.Server
	port       int
	logger     logger.Logger
}

// NewServer creates new gRPC server.
func NewServer(
	port int,
	userHandler *handler.UserServiceServer,
	jwtClient jwt.TokenClient,
	logger logger.Logger,
) *Server {
	// TODO: Create unary interceptor chain
	// - LoggingInterceptor
	// - RecoveryInterceptor
	// - AuthInterceptor (for protected methods)

	// TODO: Create grpc.Server with interceptors

	// TODO: Register UserService server to grpc.Server

	return &Server{
		grpcServer: nil, // TODO: assign created grpc.Server
		port:       port,
		logger:     logger,
	}
}

// Start starts listening on configured port.
func (s *Server) Start() error {
	// TODO: Create listener on s.port

	// TODO: Log server start on port

	// TODO: Start serving with s.grpcServer.Serve()

	return nil
}

// Stop gracefully shuts down server.
func (s *Server) Stop(ctx context.Context) error {
	// TODO: Log server shutdown

	// TODO: Use s.grpcServer.GracefulStop() or Stop() with context

	return nil
}
