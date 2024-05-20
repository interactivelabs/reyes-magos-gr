package handlers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	my_codes "reyes-magos-gr/views/volunteer/my_codes"

	"github.com/labstack/echo/v4"
)

type MyCodesHandler struct {
	VolunteersService services.VolunteersService
}

func (h MyCodesHandler) MyCodesViewHandler(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	codes, err := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, my_codes.MyCodes(codes))
}
