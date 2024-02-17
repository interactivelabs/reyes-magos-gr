package api

import (
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"strconv"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type CreateVolunteerRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address" validate:"required"`
	Address2 string `json:"address2"`
	Country  string `json:"country" validate:"required"`
	State    string `json:"state" validate:"required"`
	City     string `json:"city" validate:"required"`
	Province string `json:"province"`
	ZipCode  string `json:"zip_code" validate:"required"`
	Secret   string `json:"secret" validate:"required"`
	Passcode string `json:"passcode" validate:"required,number"`
}

type VolunteerHandler struct {
	VolunteersRepository repository.VolunteersRepository
}

func (h VolunteerHandler) CreateVolunteerApiHandler(ctx echo.Context) error {
	tr := new(CreateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return err
	}

	var volunteer model.Volunteer
	err := dto.Map(&volunteer, tr)
	if err != nil {
		return err
	}

	volunteerID, err := h.VolunteersRepository.CreateVolunteer(volunteer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, volunteerID)
}

type UpdateVolunteerRequest struct {
	VolunteerID int64  `json:"volunteer_id" validate:"required"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	Country     string `json:"country"`
	State       string `json:"state"`
	City        string `json:"city"`
	Province    string `json:"province"`
	ZipCode     string `json:"zip_code"`
	Secret      string `json:"secret"`
	Passcode    string `json:"passcode" validate:"omitempty,number"`
}

func (h VolunteerHandler) UpdateVolunteerApiHandler(ctx echo.Context) error {
	tr := new(UpdateVolunteerRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return err
	}

	var volunteer model.Volunteer
	err := dto.Map(&volunteer, tr)
	if err != nil {
		return err
	}

	err = h.VolunteersRepository.UpdateVolunteer(volunteer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, volunteer.VolunteerID)
}

type DeleteVolunteerRequest struct {
}

func (h VolunteerHandler) DeleteVolunteerApiHandler(ctx echo.Context) error {

	volunteerIDStr := ctx.Param("volunteer_id")
	volunteerID, err := strconv.ParseInt(volunteerIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid volunteer ID")
	}

	err = h.VolunteersRepository.DeleteVolunteer(volunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, volunteerID)
}
