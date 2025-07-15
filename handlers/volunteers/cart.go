package volunteers

import (
	"net/http"
	"reyes-magos-gr/services"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	VolunteersService services.VolunteersService
}

func NewCartHandler(volunteersService services.VolunteersService) *CartHandler {
	return &CartHandler{
		VolunteersService: volunteersService,
	}
}

func (h *CartHandler) CartViewHandler(ctx echo.Context) error {
	profile, perr := GetProfileHandler(ctx, h.VolunteersService)
	if perr != nil && perr.Code == http.StatusUnauthorized {
		return perr
	}
	if perr != nil && perr.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	cart, err := h.VolunteersService.GetVolunteerCartByEmail(profile.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, cart)
}

type CreateCartItemRequest struct {
	ToyID int64 `json:"toy_id" validate:"required"`
}

func (h *CartHandler) CreateCartItemHandler(ctx echo.Context) error {
	// Get the toy ID from the request body
	createCartItem := new(CreateCartItemRequest)
	if err := ctx.Bind(createCartItem); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(createCartItem); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Get the volunteer profile from the context
	profile, perr := GetProfileHandler(ctx, h.VolunteersService)
	if perr != nil && perr.Code == http.StatusUnauthorized {
		return perr
	}
	if perr != nil && perr.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	item, err := h.VolunteersService.CreateVolunteerCartItem(profile.Email, createCartItem.ToyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusCreated, map[string]int64{
		"cart_id": item,
	})
}
