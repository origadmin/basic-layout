package service

import (
	"context"

	"github.com/google/wire"

	userv1 "basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"basic-layout/multiple/multiple_sample/internal/mods/user/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// UserService is a user service.
type UserService struct {
	userv1.UnimplementedUserAPIServer
	uc  *biz.UserUsecase
	log *log.Helper
}

// NewUserService new a user service.
func NewUserService(uc *biz.UserUsecase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(log.With(logger, "module", "service/user"))}
}

// CreateUser implements user.UserAPIServer.
func (s *UserService) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	// DTO to DO conversion
	userDO := &biz.User{
		Username: req.GetUsername(),
		Nickname: req.GetNickname(),
	}

	// Call business logic
	createdUser, err := s.uc.Create(ctx, userDO)
	if err != nil {
		return nil, err
	}

	// DO to DTO conversion
	return &userv1.CreateUserResponse{
		User: toUserDTO(createdUser),
	}, nil
}

// GetUser implements user.UserAPIServer.
func (s *UserService) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	// Call business logic
	foundUser, err := s.uc.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// DO to DTO conversion
	return &userv1.GetUserResponse{
		User: toUserDTO(foundUser),
	}, nil
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewUserService, // Provides *UserService
)
