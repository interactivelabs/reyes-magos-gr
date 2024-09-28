package admin

import (
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	codes "reyes-magos-gr/views/admin/codes"

	"github.com/labstack/echo/v4"
)

type CodesHandler struct {
	CodesRepository          repository.CodesRepository
	VolunteersRepository     repository.VolunteersRepository
	VolunteerCodesRepository repository.VolunteerCodesRepository
	CodesService             services.CodesService
}

func (h CodesHandler) CodesViewHandler(ctx echo.Context) error {
	activeCodes, err := h.CodesRepository.GetUnassignedCodes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	activeVolunteers, err := h.VolunteersRepository.GetActiveVolunteers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	allVolunteersCodes, err := h.VolunteerCodesRepository.GetAllVolunteersCodes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, codes.Codes(activeCodes, allVolunteersCodes, activeVolunteers))
}

type AssignCodesRequest struct {
	VolunteerID int64   `form:"volunteer_id" validate:"required"`
	CodeIDs     []int64 `form:"code_ids" validate:"required"`
}

func (h CodesHandler) AssignCodesHandler(ctx echo.Context) error {
	acr := new(AssignCodesRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, codeID := range acr.CodeIDs {
		volunteerCode := model.VolunteerCode{
			VolunteerID: acr.VolunteerID,
			CodeID:      codeID,
		}
		_, err := h.VolunteerCodesRepository.CreateVolunteerCode(volunteerCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.Redirect(303, "/admin/codes")
}

type RemoveCodesRequest struct {
	VolunteerCodeIDs []int64 `form:"volunteer_code_ids" validate:"required"`
	CodeIDs          []int64 `form:"code_ids" validate:"required"`
}

func (h CodesHandler) RemoveCodesHandler(ctx echo.Context) error {
	acr := new(RemoveCodesRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, volunteerCodeID := range acr.VolunteerCodeIDs {

		err := h.VolunteerCodesRepository.DeleteVolunteerCode(volunteerCodeID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	for _, codeID := range acr.CodeIDs {
		err := h.CodesRepository.DeleteCode(codeID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.Redirect(303, "/admin/codes")
}

type CreateCodesRequest struct {
	Count int64 `form:"count" validate:"required,min=1,max=100"`
}

func (h CodesHandler) CreateCodesHandler(ctx echo.Context) error {
	acr := new(CreateCodesRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := h.CodesService.CreateCodeBatch(acr.Count)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(303, "/admin/codes")
}
