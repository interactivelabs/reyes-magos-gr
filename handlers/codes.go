package handlers

import (
	codes "reyes-magos-gr/components/codes"
	"reyes-magos-gr/db/repository"

	"github.com/labstack/echo/v4"
)

type CodesHandler struct {
	CodesRepository repository.CodesRepository
}

func (h CodesHandler) CodesViewHandler(ctx echo.Context) error {
	activeCodes, err := h.CodesRepository.GetActiveCodes()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return render(ctx, codes.Codes(activeCodes))
}
