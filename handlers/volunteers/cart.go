package volunteers

import (
	"net/http"
	"reyes-magos-gr/services"

	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	VolunteersService services.VolunteersService
}

func (h CartHandler) CartViewHandler(ctx echo.Context) error {
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
