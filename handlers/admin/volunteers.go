package admin

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	views "reyes-magos-gr/views/admin/volunteers"

	"github.com/labstack/echo/v4"
)

type VolunteersHandler struct {
	VolunteersService services.VolunteersService
}

func (h VolunteersHandler) VolunteersViewHandler(ctx echo.Context) error {
	volunteers, err := h.VolunteersService.GetActiveVolunteersGrupedByLocation()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, views.AdminVolunteers(volunteers))
}

func (h VolunteersHandler) VolunteersCreateHandler(ctx echo.Context) error {
	return lib.Render(ctx, views.CreateVolunteerForm())
}
