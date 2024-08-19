package volunteers

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/services"
	my_codes "reyes-magos-gr/views/volunteer/my_codes"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MyCodesHandler struct {
	VolunteersService services.VolunteersService
	CodesRepository   repository.CodesRepository
}

func (h MyCodesHandler) MyCodesViewHandler(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	codes, givenCodes, err := h.VolunteersService.GetVolunteerCodesByEmail(profile.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, my_codes.MyCodes(codes, givenCodes))
}

func (h MyCodesHandler) GiveCode(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	codeIDStr := ctx.Param("code_id")
	codeID, err := strconv.ParseInt(codeIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid code ID")
	}

	code, err := h.CodesRepository.GetCodeByID(codeID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	code.Given = 1

	err = h.CodesRepository.UpdateCode(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, my_codes.MyCodeItem(code))
}
