package volunteers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/dtos"
	volunteer "reyes-magos-gr/views/volunteer"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MyCodesHandler struct {
	CodesStore        store.CodesStore
	VolunteersService services.VolunteersService
}

func NewMyCodesHandler(
	codesStore store.CodesStore,
	volunteersService services.VolunteersService,
) *MyCodesHandler {
	return &MyCodesHandler{
		CodesStore:        codesStore,
		VolunteersService: volunteersService,
	}
}

func (h *MyCodesHandler) MyCodesViewHandler(ctx echo.Context) error {
	profile, perr := GetProfileHandler(ctx, h.VolunteersService)
	if perr != nil && perr.Code == http.StatusUnauthorized {
		return perr
	}
	if perr != nil && perr.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	codes, givenCodes, err := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, volunteer.MyCodes(codes, givenCodes))
}

func (h *MyCodesHandler) GiveCode(ctx echo.Context) error {
	codeIDStr := ctx.Param("code_id")
	codeID, err := strconv.ParseInt(codeIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid code ID")
	}

	code, err := h.CodesStore.GetCodeByID(codeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	code.Given = 1

	err = h.CodesStore.UpdateCode(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, volunteer.MyCodeItem(code))
}

type GiveCodesRequest struct {
	CodesNumber int `form:"codes_number" validate:"required"`
}

func (h *MyCodesHandler) GiveCodes(ctx echo.Context) error {
	acr := new(GiveCodesRequest)
	if err := ctx.Bind(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(acr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	profile, err := GetProfileHandler(ctx, h.VolunteersService)
	if err != nil && err.Code == http.StatusUnauthorized {
		return err
	}
	if err != nil && err.Code == http.StatusForbidden {
		return ctx.Redirect(http.StatusTemporaryRedirect, "/notvolunteer")
	}

	codes, _, cerr := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if cerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, cerr.Error())
	}

	if acr.CodesNumber > len(codes) {
		return echo.NewHTTPError(http.StatusBadRequest, "Not enough codes available")
	}

	// Select the first N codes to mark as given
	codesToGive := codes[:acr.CodesNumber]

	// Mark all selected codes as given
	for i := range codesToGive {
		codesToGive[i].Given = 1
	}

	// Update all codes in a batch
	if err := h.CodesStore.UpdateCodes(codesToGive); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Fetch updated codes for rendering
	updatedCodes, updatedGivenCodes, cerr := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if cerr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, cerr.Error())
	}

	return lib.Render(ctx, volunteer.MyCodes(updatedCodes, updatedGivenCodes))
}

func GetProfileHandler(
	ctx echo.Context,
	volunteersService services.VolunteersService,
) (dtos.Profile, *echo.HTTPError) {
	profile := lib.GetCtxProfile(ctx)
	if profile.Email == "" {
		return dtos.Profile{}, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	_, err := volunteersService.GetVolunteerByEmail(profile.Email)
	if err != nil {
		return dtos.Profile{}, echo.NewHTTPError(http.StatusForbidden, "Not a volunteer")
	}

	return profile, nil
}
