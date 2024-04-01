package handlers

import (
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	orders "reyes-magos-gr/views/orders"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersRepository repository.OrdersRepository
}

func (h OrdersHandler) CreateOrderViewHandler(ctx echo.Context) error {
	// toys, err := h.OrdersRepository.GetOrders()
	// if err != nil {
	// 	return echo.NewHTTPError(500, err.Error())
	// }

	return lib.Render(ctx, orders.CreateOrder())
}
