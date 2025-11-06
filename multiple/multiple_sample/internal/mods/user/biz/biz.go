package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is user module's biz providers.
var ProviderSet = wire.NewSet(NewUserUsecase)

// UserRepo defines the data access interface for the User entity.
// It is implemented by the data layer.
type UserRepo interface {
	Create(context.Context, *User) (*User, error)
	Get(context.Context, int64) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, int64) error
	List(ctx context.Context, current, pageSize int32) ([]*User, int32, error)
}

// User is a User model (Domain Object).
// It represents the user entity in the business domain.
type User struct {
	ID       int64
	Username string
	Nickname string
}

// UserUsecase is a User usecase.
// It contains the business logic for user-related operations.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase creates a new UserUsecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/user"))}
}

// Create creates a new User.
func (uc *UserUsecase) Create(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %v", u.Username)
	// Here you can add business logic, e.g., check if username exists.
	return uc.repo.Create(ctx, u)
}

// Get retrieves a User by their ID.
func (uc *UserUsecase) Get(ctx context.Context, id int64) (*User, error) {
	uc.log.WithContext(ctx).Infof("GetUser: %v", id)
	return uc.repo.Get(ctx, id)
}

// Update updates an existing User.
func (uc *UserUsecase) Update(ctx context.Context, u *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %v", u.ID)
	return uc.repo.Update(ctx, u)
}

// Delete deletes a User by their ID.
func (uc *UserUsecase) Delete(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: %v", id)
	return uc.repo.Delete(ctx, id)
}

// List lists all Users with pagination.
func (uc *UserUsecase) List(ctx context.Context, current, pageSize int32) ([]*User, int32, error) {
	uc.log.WithContext(ctx).Info("ListUsers")
	return uc.repo.List(ctx, current, pageSize)
}
