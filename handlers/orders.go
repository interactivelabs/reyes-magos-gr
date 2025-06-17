package handlers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	orders "reyes-magos-gr/views/orders"
	redeem "reyes-magos-gr/views/redeem-toy"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersService        services.OrdersService
	VolunteersRepository store.VolunteersRepository
}

type CreateOrderRequest struct {
	ToyID int64  `form:"toy_id" validate:"required"`
	Code  string `form:"code" validate:"required"`
}

func (h OrdersHandler) CreateOrderViewHandler(ctx echo.Context) error {
	acr := new(CreateOrderRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersService.CreateOrder(acr.ToyID, acr.Code)
	if err != nil {
		if lib.GetHTMLErrorCode(err) == http.StatusBadRequest || lib.GetHTMLErrorCode(err) == http.StatusConflict {
			return lib.Render(ctx, redeem.RedeemToyForm(acr.ToyID, acr.Code, "Codigo invalido"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	volunteer, err := h.VolunteersRepository.GetVolunteerByID(order.VolunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, orders.OrderCreatedSucessBanner(volunteer.Name))
}
