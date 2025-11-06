package service

import (
	"context"

	"github.com/google/wire"

	orderv1 "basic-layout/multiple/multiple_sample/api/v1/gen/go/order"
	"basic-layout/multiple/multiple_sample/internal/mods/order/biz"

	"github.com/go-kratos/kratos/v2/log"
)

// OrderService is an order service.
type OrderService struct {
	orderv1.UnimplementedOrderAPIServer
	uc  *biz.OrderUsecase
	log *log.Helper
}

func toOrderDTO(order *biz.Order) *orderv1.Order {
	return &orderv1.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Amount: order.Amount,
		Status: order.Status,
	}
}

// NewOrderService new an order service.
func NewOrderService(uc *biz.OrderUsecase, logger log.Logger) *OrderService {
	return &OrderService{uc: uc, log: log.NewHelper(log.With(logger, "module", "service/order"))}
}

// DeleteOrder deletes an order by ID.
func (s *OrderService) DeleteOrder(ctx context.Context, request *orderv1.DeleteOrderRequest) (*orderv1.DeleteOrderResponse, error) {
	err := s.uc.Delete(ctx, request.GetId())
	if err != nil {
		return nil, err
	}
	return &orderv1.DeleteOrderResponse{Success: true}, nil
}

// ListOrder lists all orders.
func (s *OrderService) ListOrder(ctx context.Context, request *orderv1.ListOrderRequest) (*orderv1.ListOrderResponse, error) {
	// Assuming biz.OrderUsecase.List takes context and returns a slice of biz.Order and total count
	// We need to update biz.OrderUsecase.List to accept pagination parameters and return total.
	ordersDO, total, err := s.uc.List(ctx, request.GetCurrent(), request.GetPageSize())
	if err != nil {
		return nil, err
	}

	var orderDTOs []*orderv1.Order
	for _, orderDO := range ordersDO {
		orderDTOs = append(orderDTOs, toOrderDTO(orderDO))
	}

	return &orderv1.ListOrderResponse{
			Data:  orderDTOs,
			Total: total,
		},
		nil
}

// UpdateOrder updates an existing order.
func (s *OrderService) UpdateOrder(ctx context.Context, request *orderv1.UpdateOrderRequest) (*orderv1.UpdateOrderResponse, error) {
	// DTO to DO conversion
	orderDO := &biz.Order{
		ID:     request.GetOrder().GetId(),
		UserID: request.GetOrder().GetUserId(),
		Amount: request.GetOrder().GetAmount(),
	}

	// Call business logic
	updatedOrder, err := s.uc.Update(ctx, orderDO)
	if err != nil {
		return nil, err
	}

	// DO to DTO conversion
	return &orderv1.UpdateOrderResponse{
			Order: toOrderDTO(updatedOrder),
		},
		nil
}

// CreateOrder implements order.OrderAPIServer.
func (s *OrderService) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	// DTO to DO conversion
	orderDO := &biz.Order{
		UserID: req.GetUserId(),
		Amount: req.GetAmount(),
	}

	// Call business logic
	createdOrder, err := s.uc.Create(ctx, orderDO)
	if err != nil {
		return nil, err
	}

	// DO to DTO conversion
	return &orderv1.CreateOrderResponse{
			Order: toOrderDTO(createdOrder),
		},
		nil
}

// GetOrder implements order.OrderAPIServer.
func (s *OrderService) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	foundOrder, err := s.uc.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &orderv1.GetOrderResponse{
			Order: toOrderDTO(foundOrder),
		},
		nil
}

// ProviderSet is order module's service providers.
var ProviderSet = wire.NewSet(
	NewOrderService,
)

var _ orderv1.OrderAPIServer = (*OrderService)(nil)
