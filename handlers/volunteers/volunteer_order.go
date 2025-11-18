package volunteers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	orders "reyes-magos-gr/views/orders"

	"github.com/labstack/echo/v4"
)

type VolunteerOrderHandler struct {
	OrdersService     services.OrdersService
	VolunteersService services.VolunteersService
}

func NewVolunteerOrderHandler(
	ordersService services.OrdersService,
	volunteersService services.VolunteersService,
) *VolunteerOrderHandler {
	return &VolunteerOrderHandler{
		OrdersService:     ordersService,
		VolunteersService: volunteersService,
	}
}

type CreateVolunteerOrderRequest struct {
	CartIDs    []int64 `form:"cart_id"  validate:"required"`
	Quantities []int   `form:"quantity" validate:"required"`
}

func (h *VolunteerOrderHandler) CreateOrderHandler(ctx echo.Context) error {
	createOrder := new(CreateVolunteerOrderRequest)
	if err := ctx.Bind(createOrder); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	profile, err := GetProfileHandler(ctx, h.VolunteersService)
	if err != nil && err.Code == http.StatusUnauthorized {
		return err
	}
	if err != nil && err.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	codes, _, cerr := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if cerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, cerr.Error())
	}

	oerr := h.OrdersService.CreateOrders(createOrder.CartIDs, createOrder.Quantities, codes)
	if oerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, oerr.Error())
	}

	return lib.Render(ctx, orders.OrderCreatedSucessBanner(profile.Nickname))
}
