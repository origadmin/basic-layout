package data

import (
	"basic-layout/multiple/multiple_sample/internal/mods/order/biz"
)

// toOrderDO converts an Order PO to a DO.
func toOrderDO(p *Order) *biz.Order {
	if p == nil {
		return nil
	}
	return &biz.Order{
		ID:     p.ID,
		UserID: p.UserID,
		Amount: p.Amount,
		Status: p.Status,
	}
}

// fromOrderDO converts an Order DO to a PO.
func fromOrderDO(d *biz.Order) *Order {
	if d == nil {
		return nil
	}
	return &Order{
		ID:     d.ID,
		UserID: d.UserID,
		Amount: d.Amount,
		Status: d.Status,
	}
}