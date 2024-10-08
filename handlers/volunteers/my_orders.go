package volunteers

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/views/components"
	my_orders "reyes-magos-gr/views/volunteer/my_orders"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type MyOrdersHandler struct {
	VolunteersService services.VolunteersService
	Ordersrepository  repository.OrdersRepository
}

func (h MyOrdersHandler) MyOrdersViewHandler(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	orders, err := h.VolunteersService.GetVolunteerOrdersByEmail(profile.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, my_orders.MyOrders(orders))
}

func (h MyOrdersHandler) MyOrdersCompleteHandler(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	orderIDStr := ctx.Param("order_id")
	orderId, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.Ordersrepository.GetOrderByID(orderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	order.Completed = 1
	if order.Completed == 1 {
		order.CompletedDate = time.Now().Format(time.RFC3339)
	}

	err = h.Ordersrepository.UpdateOrder(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, components.OrderCard(order))
}
