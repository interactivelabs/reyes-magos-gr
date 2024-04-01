package handlers

import (
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	redeem "reyes-magos-gr/views/redeem-toy"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RedeemToyHandler struct {
	ToysRepository repository.ToysRepository
}

func (h RedeemToyHandler) RedeemToyViewHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	toy, err := h.ToysRepository.GetToyByID(toyID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return lib.Render(ctx, redeem.RedeemToy(toy))
}
