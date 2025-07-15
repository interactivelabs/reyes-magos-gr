package handlers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	checkout "reyes-magos-gr/views/checkout"
	orders "reyes-magos-gr/views/orders"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	VolunteersStore store.VolunteersStore
	OrdersService   services.OrdersService
}

func NewOrdersHandler(
	volunteersStore store.VolunteersStore,
	ordersService services.OrdersService,

) *OrdersHandler {
	return &OrdersHandler{
		VolunteersStore: volunteersStore,
		OrdersService:   ordersService,
	}
}

type CreateOrderRequest struct {
	ToyID int64  `form:"toy_id" validate:"required"`
	Code  string `form:"code"   validate:"required"`
}

func (h *OrdersHandler) CreateOrderViewHandler(ctx echo.Context) error {
	acr := new(CreateOrderRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersService.CreateOrder(acr.ToyID, acr.Code)
	if err != nil {
		if lib.GetHTMLErrorCode(err) == http.StatusBadRequest ||
			lib.GetHTMLErrorCode(err) == http.StatusConflict {
			return lib.Render(ctx, checkout.RedeemToyForm(acr.ToyID, acr.Code, "Codigo invalido"))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	volunteer, err := h.VolunteersStore.GetVolunteerByID(order.VolunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, orders.OrderCreatedSucessBanner(volunteer.Name))
}
