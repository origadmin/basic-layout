/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"basic-layout/multiple/multiple_sample/internal/mods/order/biz"
)

type orderRepo struct {
	db  *Database
	log *log.Helper
}

func (or orderRepo) Create(ctx context.Context, order *biz.Order) (*biz.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (or orderRepo) Get(ctx context.Context, id int64) (*biz.Order, error) {
	//TODO implement me
	panic("implement me")
}

// NewOrderRepo .nfunc NewOrderRepo(db *Database, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		db:  db,
		log: log.NewHelper(log.With(logger, "module", "data/order")),
	}
}