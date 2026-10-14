package volunteers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	volunteer "reyes-magos-gr/views/volunteer"
	"strings"
	"unicode/utf8"

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
	CartIDs []int64  `form:"cart_id"`
	Notes   []string `form:"note"`
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

	if len(createOrder.CartIDs) == 0 || len(createOrder.CartIDs) != len(createOrder.Notes) {
		return renderCartContent(ctx, h.VolunteersService, profile.Email,
			"Tu carrito cambió. Revisa los artículos e intenta de nuevo.")
	}

	items := make([]services.CheckoutItem, len(createOrder.CartIDs))
	for i, cartID := range createOrder.CartIDs {
		note := strings.TrimSpace(createOrder.Notes[i])
		if note == "" {
			return renderCartContent(ctx, h.VolunteersService, profile.Email,
				"Cada artículo necesita una nota.")
		}
		if utf8.RuneCountInString(note) > maxCodeNoteLength {
			return renderCartContent(ctx, h.VolunteersService, profile.Email,
				"Las notas no pueden tener más de 200 caracteres.")
		}
		items[i] = services.CheckoutItem{CartID: cartID, Note: note}
	}

	volunteerModel, verr := h.VolunteersService.GetVolunteerByEmail(profile.Email)
	if verr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, verr.Error())
	}

	codes, _, cerr := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if cerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, cerr.Error())
	}

	oerr := h.OrdersService.CreateOrders(volunteerModel.VolunteerID, items, codes)
	if oerr != nil {
		if msg := cartErrorMessage(oerr); msg != "" {
			return renderCartContent(ctx, h.VolunteersService, profile.Email, msg)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, oerr.Error())
	}

	return lib.Render(ctx, volunteer.CartCheckoutSuccess(len(items)))
}
