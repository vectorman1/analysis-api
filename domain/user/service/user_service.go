package service

import (
	"context"
	"time"

	"github.com/vectorman1/analysis/analysis-api/domain/user/repo"

	"github.com/vectorman1/analysis/analysis-api/domain/user/model"

	"github.com/jackc/pgtype"

	validation "github.com/vectorman1/analysis/analysis-api/common/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofrs/uuid"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
)

type UserServiceContract interface {
	Login(context.Context, *user_service.LoginRequest) (*user_service.LoginResponse, error)
	Register(context.Context, *user_service.RegisterRequest) (*user_service.RegisterResponse, error)
	GetPaged(context.Context, *user_service.GetPagedRequest) (*user_service.GetPagedResponse, error)
	Create(context.Context, *user_service.CreateRequest) (*user_service.CreateResponse, error)
	Update(context.Context, *user_service.UpdateRequest) (*user_service.UpdateResponse, error)
	Delete(context.Context, *user_service.DeleteRequest) (*user_service.DeleteResponse, error)
}

type UserService struct {
	userRepository *repo.UserRepository
	config         *common.Config
}

func NewUserService(userRepository *repo.UserRepository, config *common.Config) *UserService {
	return &UserService{
		userRepository: userRepository,
		config:         config,
	}
}

// Login attempts to find a user with a matching username, verifies the password
// and returns a response with a Token or an error
func (s *UserService) Login(ctx context.Context, request *user_service.LoginRequest) (*user_service.LoginResponse, error) {
	if request.Username == "" || request.Password == "" {
		return nil, status.Error(codes.InvalidArgument, validation.WrongUsernameOrPassword)
	}

	user, err := s.userRepository.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, validation.WrongUsernameOrPassword)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, validation.WrongUsernameOrPassword)
	}

	var u string
	user.Uuid.AssignTo(&u)

	expTime := time.Now().Add(24 * 60 * time.Minute)
	claims := &common.Claims{
		Uuid:        u,
		PrivateRole: user.PrivateRole,
		StandardClaims: jwt.StandardClaims{
			Audience:  "analysis-web",
			ExpiresAt: expTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "analysis-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JwtSigningSecret))
	if err != nil {
		return nil, err
	}

	return &user_service.LoginResponse{Token: tokenString}, nil
}

// Register attempts to create a User with the corresponding
// username and hashing the password.
func (s *UserService) Register(ctx context.Context, request *user_service.RegisterRequest) (*user_service.RegisterResponse, error) {
	if request.Username == "" || request.Password == "" {
		return nil, status.Error(codes.InvalidArgument, validation.InvalidUsernameOrPassword)
	}
	if len(request.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Minimum password length is 8.")
	}
	if len(request.Username) < 6 {
		return nil, status.Error(codes.InvalidArgument, "Minimum username length is 6.")
	}

	user := &model.User{
		Uuid:        pgtype.UUID{Status: pgtype.Present},
		PrivateRole: model.Default,
		Username:    request.Username,
		Password:    request.Password,
	}
	u, _ := uuid.NewV4()
	_ = user.Uuid.Set(u)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	loginResponse, err := s.Login(
		ctx,
		&user_service.LoginRequest{
			Username: request.Username, Password: request.Password,
		})
	if err != nil {
		return nil, err
	}

	return &user_service.RegisterResponse{Token: loginResponse.Token}, nil
}

func (s *UserService) GetPaged(ctx context.Context, request *user_service.GetPagedRequest) (*user_service.GetPagedResponse, error) {
	users, total, err := s.userRepository.GetPaged(ctx, request.Filter)
	if err != nil {
		return nil, err
	}

	var protoUsers []*user_service.User
	for _, u := range *users {
		protoUsers = append(protoUsers, u.ToProto())
	}

	return &user_service.GetPagedResponse{
		Items:      protoUsers,
		TotalItems: uint64(total),
	}, nil
}

func (s *UserService) Create(ctx context.Context, request *user_service.CreateRequest) (*user_service.CreateResponse, error) {
	password := common.RandomStringWithLength(10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		PrivateRole: model.PrivateRole(request.PrivateRole),
		Username:    request.Username,
		Password:    string(hashedPassword),
	}
	u, _ := uuid.NewV4()
	user.Uuid.Set(u.Bytes())

	err = s.userRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user_service.CreateResponse{Password: password}, nil
}

func (s *UserService) Update(ctx context.Context, request *user_service.UpdateRequest) (*user_service.UpdateResponse, error) {
	panic("implement me")
}

func (s *UserService) Delete(ctx context.Context, request *user_service.DeleteRequest) (*user_service.DeleteResponse, error) {
	panic("implement me")
}
