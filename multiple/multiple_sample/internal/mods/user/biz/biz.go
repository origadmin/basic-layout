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
	// Define other methods like Update, Delete, List etc.
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