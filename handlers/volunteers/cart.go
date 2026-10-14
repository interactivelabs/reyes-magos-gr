package volunteers

import (
	"errors"
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	volunteer "reyes-magos-gr/views/volunteer"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	CartsStore        store.CartsStore
	VolunteersService services.VolunteersService
}

func NewCartHandler(
	cartsStore store.CartsStore, volunteersService services.VolunteersService) *CartHandler {
	return &CartHandler{
		CartsStore:        cartsStore,
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

	codes, _, cerr := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if cerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, cerr.Error())
	}

	return lib.Render(ctx, volunteer.Cart(cart, len(codes)))
}

type CreateCartItemRequest struct {
	ToyID int64 `form:"toy_id" validate:"required"`
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

	_, err := h.VolunteersService.CreateVolunteerCartItem(profile.Email, createCartItem.ToyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, volunteer.CartItemCreated())
}

func (h *CartHandler) DuplicateCartItemHandler(ctx echo.Context) error {
	return h.mutateCartItem(ctx, h.VolunteersService.DuplicateVolunteerCartItem)
}

func (h *CartHandler) DeleteCartItemHandler(ctx echo.Context) error {
	return h.mutateCartItem(ctx, h.VolunteersService.DeleteVolunteerCartItem)
}

// mutateCartItem runs a cart row action and re-renders the whole cart, so the
// summary counts and the add buttons' enabled state stay in sync.
func (h *CartHandler) mutateCartItem(
	ctx echo.Context,
	mutate func(email string, cartID int64) error,
) error {
	cartID, err := strconv.ParseInt(ctx.Param("cart_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Cart ID")
	}

	profile, perr := GetProfileHandler(ctx, h.VolunteersService)
	if perr != nil && perr.Code == http.StatusUnauthorized {
		return perr
	}
	if perr != nil && perr.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	errorMsg := ""
	if err := mutate(profile.Email, cartID); err != nil {
		errorMsg = cartErrorMessage(err)
		if errorMsg == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return renderCartContent(ctx, h.VolunteersService, profile.Email, errorMsg)
}

func renderCartContent(
	ctx echo.Context,
	volunteersService services.VolunteersService,
	email string,
	errorMsg string,
) error {
	cart, err := volunteersService.GetVolunteerCartByEmail(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	codes, _, err := volunteersService.GetVolunteerCodesByEmail(email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, volunteer.CartContent(cart, len(codes), errorMsg))
}

// cartErrorMessage maps known cart errors to a message for the volunteer;
// it returns "" for unexpected errors.
func cartErrorMessage(err error) string {
	switch {
	case errors.Is(err, services.ErrNotEnoughCodes):
		return "No tienes suficientes códigos disponibles para estos artículos."
	case errors.Is(err, services.ErrCartItemNotFound):
		return "Uno de los artículos ya no está en tu carrito. Revisa tu carrito e intenta de nuevo."
	default:
		return ""
	}
}
