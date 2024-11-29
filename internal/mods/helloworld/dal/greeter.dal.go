/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package dal

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"origadmin/basic-layout/internal/mods/helloworld/dto"
)

type greeterRepo struct {
	db  *Database
	log *log.Helper
}

func (g greeterRepo) Save(ctx context.Context, greeter *dto.Greeter) (*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterRepo) Update(ctx context.Context, greeter *dto.Greeter) (*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterRepo) FindByID(ctx context.Context, s string, param *dto.GreeterQueryParam) (*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterRepo) ListByHello(ctx context.Context, s string, param *dto.GreeterQueryParam) ([]*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterRepo) ListAll(ctx context.Context, param *dto.GreeterQueryParam) ([]*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterDal .
func NewGreeterDal(db *Database, logger log.Logger) dto.GreeterRepo {
	return &greeterRepo{
		db:  db,
		log: log.NewHelper(logger),
	}
}
