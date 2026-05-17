// internal/transport/grpc/handler/user_handler.go
package handler

import (
	"context"

	"github.com/artyomstank/virtual_deanery/internal/domain/service"
	"github.com/artyomstank/virtual_deanery/pkg/logger"
)

// UserServiceServer implements gRPC UserService.
type UserServiceServer struct {
	svc    service.UserService
	logger logger.Logger
	// TODO: Add generated proto server interface
}

// NewUserServiceServer creates new gRPC user service server.
func NewUserServiceServer(svc service.UserService, logger logger.Logger) *UserServiceServer {
	return &UserServiceServer{
		svc:    svc,
		logger: logger,
	}
}

// GetUser implements GetUser RPC.
func (s *UserServiceServer) GetUser(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: Parse request (proto message)

	// TODO: Extract user ID from request or context

	// TODO: Call s.svc.GetUser()

	// TODO: Map response to proto message

	// TODO: Handle errors and map to gRPC codes

	return nil, nil
}

// RegisterUser implements RegisterUser RPC.
func (s *UserServiceServer) RegisterUser(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: Parse request proto

	// TODO: Call s.svc.RegisterUser()

	// TODO: Map response to proto

	// TODO: Handle errors and map to gRPC codes

	return nil, nil
}

// LoginUser implements LoginUser RPC.
func (s *UserServiceServer) LoginUser(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: Parse request proto

	// TODO: Call s.svc.LoginUser()

	// TODO: Return tokens in response

	// TODO: Handle errors

	return nil, nil
}

// UpdateUser implements UpdateUser RPC.
func (s *UserServiceServer) UpdateUser(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: Parse request proto

	// TODO: Extract user ID from context

	// TODO: Call s.svc.UpdateUser()

	// TODO: Map response to proto

	// TODO: Handle errors

	return nil, nil
}

// DeleteUser implements DeleteUser RPC.
func (s *UserServiceServer) DeleteUser(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: Parse request proto

	// TODO: Extract user ID from context

	// TODO: Call s.svc.DeleteUser()

	// TODO: Handle errors and return Empty response

	return nil, nil
}

// ListUsers implements ListUsers RPC.
func (s *UserServiceServer) ListUsers(ctx context.Context, req interface{}) (interface{}, error) {
	// TODO: Parse request proto (pagination)

	// TODO: Call s.svc.ListUsers()

	// TODO: Map response to proto list

	// TODO: Handle errors

	return nil, nil
}
