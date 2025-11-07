package service

import (
	"context"

	"github.com/google/wire"

	userv1 "basic-layout/multiple/multiple_sample/api/v1/gen/go/user"
	"basic-layout/multiple/multiple_sample/internal/features/user/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// UserService is a user service.
type UserService struct {
	userv1.UnimplementedUserAPIServer
	uc  *biz.UserUsecase
	log *log.Helper
}

// toUserDTO converts a biz.User domain object to a userv1.User DTO.
func toUserDTO(do *biz.User) *userv1.User {
	if do == nil {
		return nil
	}
	return &userv1.User{
		Id:       do.ID,
		Username: do.Username,
		Nickname: do.Nickname,
	}
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

// UpdateUser implements user.UserAPIServer.
func (s *UserService) UpdateUser(ctx context.Context, request *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	userDO := &biz.User{
		ID:       request.GetUser().GetId(),
		Username: request.GetUser().GetUsername(),
		Nickname: request.GetUser().GetNickname(),
	}

	updatedUser, err := s.uc.Update(ctx, userDO)
	if err != nil {
		return nil, err
	}

	return &userv1.UpdateUserResponse{
		User: toUserDTO(updatedUser),
	}, nil
}

// DeleteUser implements user.UserAPIServer.
func (s *UserService) DeleteUser(ctx context.Context, request *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	err := s.uc.Delete(ctx, request.GetId())
	if err != nil {
		return nil, err
	}
	return &userv1.DeleteUserResponse{Success: true}, nil
}

// ListUser implements user.UserAPIServer.
func (s *UserService) ListUser(ctx context.Context, request *userv1.ListUserRequest) (*userv1.ListUserResponse, error) {
	usersDO, total, err := s.uc.List(ctx, request.GetCurrent(), request.GetPageSize())
	if err != nil {
		return nil, err
	}

	var userDTOs []*userv1.User
	for _, userDO := range usersDO {
		userDTOs = append(userDTOs, toUserDTO(userDO))
	}

	return &userv1.ListUserResponse{
			Data:  userDTOs,
			Total: total,
		},
		nil
}

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	NewUserService, // Provides *UserService
)
