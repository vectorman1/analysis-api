package service

import (
	"time"

	"github.com/vectorman1/analysis/analysis-api/model"

	dbmodel "github.com/vectorman1/analysis/analysis-api/model/db"

	"golang.org/x/crypto/bcrypt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/pgtype"

	"github.com/vectorman1/analysis/analysis-api/common"

	"github.com/dgrijalva/jwt-go"
	"github.com/vectorman1/analysis/analysis-api/db"
	"github.com/vectorman1/analysis/analysis-api/generated/user_service"
)

type userService interface {
	Login(*user_service.LoginRequest) (*user_service.LoginResponse, error)
	Register(*user_service.RegisterRequest) (*user_service.RegisterResponse, error)
}

type UserService struct {
	userService
	userRepository *db.UserRepository
	config         *common.Config
}

func NewUserService(userRepository *db.UserRepository, config *common.Config) *UserService {
	return &UserService{
		userRepository: userRepository,
		config:         config,
	}
}

// Login attempts to find a user with a matching username and hashed password
// and returns a response with a Token or an error
func (s *UserService) Login(request *user_service.LoginRequest) (*user_service.LoginResponse, error) {
	user, err := s.userRepository.Get(request.Username, request.Password)
	if err != nil {
		return nil, err
	}
	var u string
	user.Uuid.AssignTo(&u)

	expTime := time.Now().Add(24 * 60 * time.Minute)
	claims := &model.Claims{
		Uuid: u,
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
func (s *UserService) Register(request *user_service.RegisterRequest) (*user_service.RegisterResponse, error) {
	user := &dbmodel.User{
		Uuid:        pgtype.UUID{Status: pgtype.Present},
		PrivateRole: dbmodel.Default,
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
	user.CreatedAt = pgtype.Timestamptz{Time: time.Now(), Status: pgtype.Present}
	user.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Status: pgtype.Present}

	err = s.userRepository.Create(user)
	if err != nil {
		return nil, err
	}

	loginResponse, err := s.Login(
		&user_service.LoginRequest{
			Username: request.Username, Password: request.Password,
		})
	if err != nil {
		return nil, err
	}

	return &user_service.RegisterResponse{Token: loginResponse.Token}, nil
}
