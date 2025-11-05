package service

import (
	"basic-layout/multiple/multiple_sample/api/v1/mods/order"
	"basic-layout/multiple/multiple_sample/internal/mods/order/biz"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

// OrderService is an order service.
type OrderService struct {
	order.UnimplementedOrderAPIServer
	uc  *biz.OrderUsecase
	log *log.Helper
}

// NewOrderService new an order service.
func NewOrderService(uc *biz.OrderUsecase, logger log.Logger) *OrderService {
	return &OrderService{uc: uc, log: log.NewHelper(log.With(logger, "module", "service/order"))}
}

// CreateOrder implements order.OrderAPIServer.
func (s *OrderService) CreateOrder(ctx context.Context, req *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	// DTO to DO conversion
	orderDO := &biz.Order{
		UserID: req.GetUserId(), // Corrected method name
		Amount: req.GetAmount(),
	}

	// Call business logic
	createdOrder, err := s.uc.Create(ctx, orderDO)
	if err != nil {
		return nil, err
	}

	// DO to DTO conversion
	return &order.CreateOrderResponse{
		Order: toOrderDTO(createdOrder),
	}, nil
}

// GetOrder implements order.OrderAPIServer.
func (s *OrderService) GetOrder(ctx context.Context, req *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	foundOrder, err := s.uc.Get(ctx, req.GetId()) // Corrected method name
	if err != nil {
		return nil, err
	}

	return &order.GetOrderResponse{
		Order: toOrderDTO(foundOrder),
	}, nil
}

// ... other methods like UpdateOrder, DeleteOrder, ListOrder can be implemented here
