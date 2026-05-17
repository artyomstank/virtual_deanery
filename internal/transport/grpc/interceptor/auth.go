// internal/transport/grpc/interceptor/auth.go
package interceptor

import (
	"context"

	"github.com/artyomstank/virtual_deanery/pkg/jwt"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
	"google.golang.org/grpc"
)

// AuthInterceptor validates JWT token from metadata.
func AuthInterceptor(jwtClient jwt.TokenClient, logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Extract token from metadata (authorization header)
		// Format: "authorization: Bearer <token>"

		// TODO: Validate token using jwtClient.ValidateToken()

		// TODO: Extract claims from token

		// TODO: Put user_id in context

		// TODO: Handle invalid/expired token → return Unauthenticated status

		return handler(ctx, req)
	}
}

// LoggingInterceptor logs gRPC requests.
func LoggingInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Log request: method, request info
		resp, err := handler(ctx, req)
		// TODO: Log response: error status if any
		return resp, err
	}
}

// RecoveryInterceptor recovers from panics.
func RecoveryInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: Implement panic recovery
		// TODO: Return Internal status on panic
		return handler(ctx, req)
	}
}
