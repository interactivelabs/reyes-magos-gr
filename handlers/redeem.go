package handlers

import (
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	redeem "reyes-magos-gr/views/redeem"

	"github.com/labstack/echo/v4"
)

type RedeemHandler struct {
	ToysRepository repository.ToysRepository
}

func (h RedeemHandler) RedeemViewHandler(ctx echo.Context) error {
	toys, err := h.ToysRepository.GetToys()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return lib.Render(ctx, redeem.Redeem(toys))
}
