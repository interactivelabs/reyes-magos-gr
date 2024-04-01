package handlers

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	orders "reyes-magos-gr/views/orders"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersRepository repository.OrdersRepository
}

type CreateOrderRequest struct {
	ToyID int64 `form:"toy_id" validate:"required"`
	Code  int64 `form:"code" validate:"required"`
}

func (h OrdersHandler) CreateOrderViewHandler(ctx echo.Context) error {
	acr := new(CreateOrderRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return err
	}

	// toys, err := h.OrdersRepository.GetOrders()
	// if err != nil {
	// 	return echo.NewHTTPError(500, err.Error())
	// }

	return lib.Render(ctx, orders.CreateOrder())
}
