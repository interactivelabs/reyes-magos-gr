package admin

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	views "reyes-magos-gr/views/admin/volunteers"
	"strconv"

	"github.com/labstack/echo/v4"
)

type VolunteersHandler struct {
	VolunteersService    services.VolunteersService
	VolunteersRepository store.VolunteersRepository
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

func (h VolunteersHandler) VolunteersCreatePostHandler(ctx echo.Context) error {
	tr := new(services.CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err := h.VolunteersService.CreateAndGetVolunteer(*tr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, views.NewVolunteerRow(volunteer))
}

func (h VolunteersHandler) VolunteersUpdateViewHandler(ctx echo.Context) error {
	volunteerID, err := strconv.ParseInt(ctx.Param("volunteer_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err := h.VolunteersRepository.GetVolunteerByID(volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, views.UpdateVolunteerForm(volunteer))
}

func (h VolunteersHandler) VolunteersUpdatePutHandler(ctx echo.Context) error {
	volunteerIDStr := ctx.Param("volunteer_id")
	volunteerID, err := strconv.ParseInt(volunteerIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid volunteer ID")
	}

	tr := new(services.CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err := h.VolunteersService.UpdateVolunteer(*tr, volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, views.VolunteerRow(volunteer))
}

func (h VolunteersHandler) VolunteersDeleteHandler(ctx echo.Context) error {
	volunteerIDStr := ctx.Param("volunteer_id")
	volunteerID, err := strconv.ParseInt(volunteerIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid volunteer ID")
	}

	err = h.VolunteersRepository.DeleteVolunteer(volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}
