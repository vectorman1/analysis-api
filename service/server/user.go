package server

import (
	"context"

	"google.golang.org/grpc/status"

	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
	"github.com/vectorman1/analysis/analysis-api/service"
)

type UserServiceServer struct {
	userService *service.UserService
	user_service.UnimplementedUserServiceServer
}

func NewUserServiceServer(userService *service.UserService) *UserServiceServer {
	return &UserServiceServer{userService: userService}
}

func (s *UserServiceServer) Login(ctx context.Context, req *user_service.LoginRequest) (*user_service.LoginResponse, error) {
	result, err := s.userService.Login(req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, st.Err()
		}
		return nil, err
	}

	return result, nil
}
func (s *UserServiceServer) Register(ctx context.Context, req *user_service.RegisterRequest) (*user_service.RegisterResponse, error) {
	result, err := s.userService.Register(req)
	if err != nil {
		return nil, err
	}

	return result, nil
}
