package handlers

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	redeem "reyes-magos-gr/views/redeem-toy"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RedeemToyHandler struct {
	ToysRepository repository.ToysRepository
}

type RedeemToyViewRequest struct {
	Code string `query:"code"`
}

func (h RedeemToyHandler) RedeemToyViewHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	toy, err := h.ToysRepository.GetToyByID(toyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cr := new(CatalogHandlerRequest)
	if err := ctx.Bind(cr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	code := cr.Code

	return lib.Render(ctx, redeem.RedeemToy(toy, code))
}
