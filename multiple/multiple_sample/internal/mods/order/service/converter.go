package service

import (
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/order"
	"basic-layout/multiple/multiple_sample/internal/mods/order/biz"
)

// toOrderDTO converts an Order DO to a DTO.
func toOrderDTO(d *biz.Order) *order.Order {
	if d == nil {
		return nil
	}
	return &order.Order{
		Id:     d.ID,
		UserId: d.UserID,
		Amount: d.Amount,
		Status: d.Status,
	}
}