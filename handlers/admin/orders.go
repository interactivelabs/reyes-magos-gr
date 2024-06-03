package handlers_admin

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	ordersView "reyes-magos-gr/views/admin/orders"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersRepository repository.OrdersRepository
}

func (h OrdersHandler) OrdersViewHandler(ctx echo.Context) error {
	orders, err := h.OrdersRepository.GetAllActiveOrders()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.Orders(orders))
}

func (h OrdersHandler) UpdateOrderViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, ordersView.Create())
}
