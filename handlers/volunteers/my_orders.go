package volunteers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	"reyes-magos-gr/views/components"
	volunteer "reyes-magos-gr/views/volunteer"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type MyOrdersHandler struct {
	VolunteersService services.VolunteersService
	Ordersrepository  store.OrdersStore
}

func (h MyOrdersHandler) MyOrdersViewHandler(ctx echo.Context) error {
	profile, perr := GetProfileHandler(ctx, h.VolunteersService)
	if perr != nil && perr.Code == http.StatusUnauthorized {
		return perr
	}
	if perr != nil && perr.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	orders, err := h.VolunteersService.GetVolunteerOrdersByEmail(profile.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, volunteer.MyOrders(orders))
}

func (h MyOrdersHandler) MyOrdersCompleteHandler(ctx echo.Context) error {
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
