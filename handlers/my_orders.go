package handlers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	my_orders "reyes-magos-gr/views/volunteer/my_orders"

	"github.com/labstack/echo/v4"
)

type MyOrdersHandler struct {
	VolunteersService services.VolunteersService
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
