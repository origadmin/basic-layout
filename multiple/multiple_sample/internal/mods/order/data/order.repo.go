/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package data

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"

	"basic-layout/multiple/multiple_sample/internal/mods/order/biz"
)

type orderRepo struct {
	db  *Data
	log *log.Helper
}

// Create creates a new order in the data store.
func (or *orderRepo) Create(ctx context.Context, order *biz.Order) (*biz.Order, error) {
	or.log.WithContext(ctx).Infof("data: createOrder: %v", order.UserID)
	// Simulate database insertion
	order.ID = 1 // Assign a dummy ID
	return order, nil
}

// Get retrieves an order by its ID from the data store.
func (or *orderRepo) Get(ctx context.Context, id int64) (*biz.Order, error) {
	or.log.WithContext(ctx).Infof("data: getOrder: %v", id)
	// Simulate database retrieval
	if id == 1 {
		return &biz.Order{ID: 1, UserID: 101, Amount: 99.99, Status: "completed"}, nil
	}
	return nil, errors.New("order not found")
}

// Update updates an existing order in the data store.
func (or *orderRepo) Update(ctx context.Context, order *biz.Order) (*biz.Order, error) {
	or.log.WithContext(ctx).Infof("data: updateOrder: %v", order.ID)
	// Simulate database update
	if order.ID == 1 {
		return order, nil
	}
	return nil, errors.New("order not found for update")
}

// Delete deletes an order by its ID from the data store.
func (or *orderRepo) Delete(ctx context.Context, id int64) error {
	or.log.WithContext(ctx).Infof("data: deleteOrder: %v", id)
	// Simulate database deletion
	if id == 1 {
		return nil
	}
	return errors.New("order not found for deletion")
}

// List retrieves a list of orders with pagination from the data store.
func (or *orderRepo) List(ctx context.Context, current, pageSize int32) ([]*biz.Order, int32, error) {
	or.log.WithContext(ctx).Infof("data: listOrders, current: %d, pageSize: %d", current, pageSize)
	// Simulate database listing with pagination
	dummyOrders := []*biz.Order{
		{ID: 1, UserID: 101, Amount: 99.99, Status: "completed"},
		{ID: 2, UserID: 102, Amount: 150.00, Status: "pending"},
	}
	total := int32(len(dummyOrders))

	start := (current - 1) * pageSize
	end := start + pageSize

	if start < 0 || start >= total {
		return []*biz.Order{}, total, nil
	}
	if end > total {
		end = total
	}

	return dummyOrders[start:end], total, nil
}

// NewOrderRepo creates a new OrderRepo.
func NewOrderRepo(db *Data, logger log.Logger) biz.OrderRepo {
	return &orderRepo{
		db:  db,
		log: log.NewHelper(log.With(logger, "module", "data/order")),
	}
}
