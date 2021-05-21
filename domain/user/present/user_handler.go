package present

import (
	"context"

	service2 "github.com/vectorman1/analysis/analysis-api/domain/user/service"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
)

type UserServiceServer struct {
	userService *service2.UserService
	user_service.UnimplementedUserServiceServer
}

func NewUserServiceServer(userService *service2.UserService) *UserServiceServer {
	return &UserServiceServer{userService: userService}
}

func (s *UserServiceServer) Login(ctx context.Context, req *user_service.LoginRequest) (*user_service.LoginResponse, error) {
	result, err := s.userService.Login(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return result, nil
}
func (s *UserServiceServer) Register(ctx context.Context, req *user_service.RegisterRequest) (*user_service.RegisterResponse, error) {
	result, err := s.userService.Register(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return result, nil
}

func (s *UserServiceServer) GetPaged(ctx context.Context, req *user_service.GetPagedRequest) (*user_service.GetPagedResponse, error) {
	result, err := s.userService.GetPaged(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return result, nil
}

func (s *UserServiceServer) Create(ctx context.Context, req *user_service.CreateRequest) (*user_service.CreateResponse, error) {
	result, err := s.userService.Create(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return result, nil
}
func (s *UserServiceServer) Update(ctx context.Context, req *user_service.UpdateRequest) (*user_service.UpdateResponse, error) {
	result, err := s.userService.Update(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return result, nil
}

func (s *UserServiceServer) Delete(ctx context.Context, req *user_service.DeleteRequest) (*user_service.DeleteResponse, error) {
	result, err := s.userService.Delete(ctx, req)
	if err != nil {
		return nil, common.GetErrorStatus(err)
	}

	return result, nil
}
