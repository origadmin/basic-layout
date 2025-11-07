/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package data

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"

	"basic-layout/multiple/multiple_sample/internal/features/user/biz"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// Create creates a new user in the data store.
func (ur *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	ur.log.WithContext(ctx).Infof("data: createUser: %v", user.Username)
	// Simulate database insertion
	user.ID = 1 // Assign a dummy ID
	return user, nil
}

// Get retrieves a user by their ID from the data store.
func (ur *userRepo) Get(ctx context.Context, id int64) (*biz.User, error) {
	ur.log.WithContext(ctx).Infof("data: getUser: %v", id)
	// Simulate database retrieval
	if id == 1 {
		return &biz.User{ID: 1, Username: "testuser", Nickname: "Test User"}, nil
	}
	return nil, errors.New("user not found")
}

// Update updates an existing user in the data store.
func (ur *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	ur.log.WithContext(ctx).Infof("data: updateUser: %v", user.ID)
	// Simulate database update
	if user.ID == 1 {
		return user, nil
	}
	return nil, errors.New("user not found for update")
}

// Delete deletes a user by their ID from the data store.
func (ur *userRepo) Delete(ctx context.Context, id int64) error {
	ur.log.WithContext(ctx).Infof("data: deleteUser: %v", id)
	// Simulate database deletion
	if id == 1 {
		return nil
	}
	return errors.New("user not found for deletion")
}

// List retrieves a list of users with pagination from the data store.
func (ur *userRepo) List(ctx context.Context, current, pageSize int32) ([]*biz.User, int32, error) {
	ur.log.WithContext(ctx).Infof("data: listUsers, current: %d, pageSize: %d", current, pageSize)
	// Simulate database listing with pagination
	dummyUsers := []*biz.User{
		{ID: 1, Username: "testuser1", Nickname: "Test User 1"},
		{ID: 2, Username: "testuser2", Nickname: "Test User 2"},
	}
	total := int32(len(dummyUsers))

	start := (current - 1) * pageSize
	end := start + pageSize

	if start < 0 || start >= total {
		return []*biz.User{}, total, nil
	}
	if end > total {
		end = total
	}

	return dummyUsers[start:end], total, nil
}

// NewUserRepo creates a new UserRepo.
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "data/user")),
	}
}
