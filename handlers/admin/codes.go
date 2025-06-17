package admin

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	codes "reyes-magos-gr/views/admin/codes"

	"github.com/labstack/echo/v4"
)

type CodesHandler struct {
	CodesStore          store.CodesStore
	VolunteersStore     store.VolunteersStore
	VolunteerCodesStore store.VolunteerCodesStore
	CodesService        services.CodesService
}

func NewCodesHandler(
	codesStore store.CodesStore,
	volunteersStore store.VolunteersStore,
	volunteerCodesStore store.VolunteerCodesStore,
	codesService services.CodesService,
) *CodesHandler {
	return &CodesHandler{
		CodesStore:          codesStore,
		VolunteersStore:     volunteersStore,
		VolunteerCodesStore: volunteerCodesStore,
		CodesService:        codesService,
	}
}

func (h *CodesHandler) CodesViewHandler(ctx echo.Context) error {
	activeCodes, err := h.CodesStore.GetUnassignedCodes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	activeVolunteers, err := h.VolunteersStore.GetActiveVolunteers()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	allVolunteersCodes, err := h.VolunteerCodesStore.GetAllVolunteersCodes()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, codes.Codes(activeCodes, allVolunteersCodes, activeVolunteers))
}

type AssignCodesRequest struct {
	VolunteerID int64   `form:"volunteer_id" validate:"required"`
	CodeIDs     []int64 `form:"code_ids"     validate:"required"`
}

func (h *CodesHandler) AssignCodesHandler(ctx echo.Context) error {
	acr := new(AssignCodesRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, codeID := range acr.CodeIDs {
		volunteerCode := models.VolunteerCode{
			VolunteerID: acr.VolunteerID,
			CodeID:      codeID,
		}
		_, err := h.VolunteerCodesStore.CreateVolunteerCode(volunteerCode)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.Redirect(303, "/admin/codes")
}

type RemoveCodesRequest struct {
	VolunteerCodeIDs []int64 `form:"volunteer_code_ids" validate:"required"`
	CodeIDs          []int64 `form:"code_ids"           validate:"required"`
}

func (h *CodesHandler) RemoveCodesHandler(ctx echo.Context) error {
	acr := new(RemoveCodesRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	for _, volunteerCodeID := range acr.VolunteerCodeIDs {

		err := h.VolunteerCodesStore.DeleteVolunteerCode(volunteerCodeID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	for _, codeID := range acr.CodeIDs {
		err := h.CodesStore.DeleteCode(codeID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return ctx.Redirect(303, "/admin/codes")
}

type CreateCodesRequest struct {
	Count int64 `form:"count" validate:"required,min=1,max=100"`
}

func (h *CodesHandler) CreateCodesHandler(ctx echo.Context) error {
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
