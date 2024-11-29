/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package dal

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"origadmin/basic-layout/internal/mods/secondworld/dto"
)

type greeterDal struct {
	db  *Database
	log *log.Helper
}

func (g greeterDal) Save(ctx context.Context, greeter *dto.Greeter) (*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterDal) Update(ctx context.Context, greeter *dto.Greeter) (*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterDal) FindByID(ctx context.Context, s string, param *dto.GreeterQueryParam) (*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterDal) ListByHello(ctx context.Context, s string, param *dto.GreeterQueryParam) ([]*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

func (g greeterDal) ListAll(ctx context.Context, param *dto.GreeterQueryParam) ([]*dto.Greeter, error) {
	//TODO implement me
	panic("implement me")
}

// NewGreeterDal .
func NewGreeterDal(db *Database, logger log.Logger) dto.GreeterRepo {
	return &greeterDal{
		db:  db,
		log: log.NewHelper(logger),
	}
}
