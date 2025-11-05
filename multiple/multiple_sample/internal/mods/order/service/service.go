package service

import (
	"basic-layout/multiple/multiple_sample/api/v1/gen/go/order"
	"github.com/google/wire"
)

// ProviderSet is order module's service providers.
var ProviderSet = wire.NewSet(NewOrderService,
	wire.Bind(new(order.OrderAPIServer), new(*OrderService)),
)
