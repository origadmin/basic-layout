/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is order module's biz providers.
var ProviderSet = wire.NewSet(NewOrderUsecase)

// OrderRepo defines the data access interface for the Order entity.
// It is implemented by the data layer.
type OrderRepo interface {
	Create(context.Context, *Order) (*Order, error)
	Get(context.Context, int64) (*Order, error)
	Update(context.Context, *Order) (*Order, error)
	Delete(context.Context, int64) error
	List(ctx context.Context, current, pageSize int32) ([]*Order, int32, error)
}

// Order is an Order model (Domain Object).
// It represents the order entity in the business domain.
type Order struct {
	ID     int64
	UserID int64
	Amount float64
	Status string
}

// OrderUsecase is an Order usecase.
// It contains the business logic for order-related operations.
type OrderUsecase struct {
	repo OrderRepo
	log  *log.Helper
}

// NewOrderUsecase creates a new OrderUsecase.
func NewOrderUsecase(repo OrderRepo, logger log.Logger) *OrderUsecase {
	return &OrderUsecase{repo: repo, log: log.NewHelper(log.With(logger, "module", "usecase/order"))}
}

// Create creates a new Order.
func (uc *OrderUsecase) Create(ctx context.Context, o *Order) (*Order, error) {
	uc.log.WithContext(ctx).Infof("CreateOrder: %v", o.ID)
	// Here you can add business logic, e.g., check if order exists.
	return uc.repo.Create(ctx, o)
}

// Get retrieves an Order by its ID.
func (uc *OrderUsecase) Get(ctx context.Context, id int64) (*Order, error) {
	uc.log.WithContext(ctx).Infof("GetOrder: %v", id)
	return uc.repo.Get(ctx, id)
}

// Update updates an existing Order.
func (uc *OrderUsecase) Update(ctx context.Context, o *Order) (*Order, error) {
	uc.log.WithContext(ctx).Infof("UpdateOrder: %v", o.ID)
	return uc.repo.Update(ctx, o)
}

// Delete deletes an Order by its ID.
func (uc *OrderUsecase) Delete(ctx context.Context, id int64) error {
	uc.log.WithContext(ctx).Infof("DeleteOrder: %v", id)
	return uc.repo.Delete(ctx, id)
}

// List lists all Orders with pagination.
func (uc *OrderUsecase) List(ctx context.Context, current, pageSize int32) ([]*Order, int32, error) {
	uc.log.WithContext(ctx).Info("ListOrders")
	return uc.repo.List(ctx, current, pageSize)
}
