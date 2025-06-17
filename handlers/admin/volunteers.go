package admin

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	views "reyes-magos-gr/views/admin/volunteers"
	"strconv"

	"github.com/dranikpg/dto-mapper"
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

type CreateVolunteerRequest struct {
	Name     string `form:"name" validate:"required"`
	Email    string `form:"email" validate:"required"`
	Phone    string `form:"phone"`
	Address  string `form:"address" validate:"required"`
	Address2 string `form:"address2"`
	Country  string `form:"country" validate:"required"`
	State    string `form:"state" validate:"required"`
	City     string `form:"city" validate:"required"`
	Province string `form:"province"`
	ZipCode  string `form:"zip_code" validate:"required"`
}

func (h VolunteersHandler) VolunteersCreatePostHandler(ctx echo.Context) error {
	tr := new(CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var volunteer models.Volunteer
	err := dto.Map(&volunteer, tr)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err = h.VolunteersService.CreateAndGetVolunteer(volunteer)
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

	tr := new(CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var volunteer models.Volunteer
	err = dto.Map(&volunteer, tr)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	volunteer, err = h.VolunteersService.UpdateVolunteer(volunteer, volunteerID)
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
