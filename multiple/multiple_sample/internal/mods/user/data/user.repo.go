/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"basic-layout/multiple/multiple_sample/internal/mods/user/biz"
)

type userRepo struct {
	db  *Database
	log *log.Helper
}

func (ur userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	//TODO implement me
	panic("implement me")
}

func (ur userRepo) Get(ctx context.Context, id int64) (*biz.User, error) {
	//TODO implement me
	panic("implement me")
}

// NewUserRepo .
func NewUserRepo(db *Database, logger log.Logger) biz.UserRepo {
	return &userRepo{
		db:  db,
		log: log.NewHelper(log.With(logger, "module", "data/user")),
	}
}