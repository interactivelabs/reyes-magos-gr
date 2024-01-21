package handlers

import (
	"net/http"
	codes "reyes-magos-gr/components/codes"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"

	"github.com/labstack/echo/v4"
)

type CodesHandler struct {
	CodesRepository          repository.CodesRepository
	VolunteersRepository     repository.VolunteersRepository
	VolunteerCodesRepository repository.VolunteerCodesRepository
}

func (h CodesHandler) CodesViewHandler(ctx echo.Context) error {
	activeCodes, err := h.CodesRepository.GetUnassignedCodes()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	activeVolunteers, err := h.VolunteersRepository.GetActiveVolunteers()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	allVolunteersCodes, err := h.VolunteerCodesRepository.GetAllVolunteersCodes()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return render(ctx, codes.Codes(activeCodes, allVolunteersCodes, activeVolunteers))
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
		return err
	}

	for _, codeID := range acr.CodeIDs {
		volunteerCode := model.VolunteerCode{
			VolunteerID: acr.VolunteerID,
			CodeID:      codeID,
		}
		_, err := h.VolunteerCodesRepository.CreateVolunteerCode(volunteerCode)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
	}

	return ctx.Redirect(303, "/admin/codes")
}
