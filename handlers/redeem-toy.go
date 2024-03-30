package handlers

import (
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	redeem "reyes-magos-gr/views/redeem-toy"
	"strconv"

	"github.com/labstack/echo/v4"
)

// RedeemToyHandler handles the redemption of toys.
type RedeemToyHandler struct {
	ToysRepository repository.ToysRepository
}

// RedeemToyViewHandler handles the HTTP request to redeem a toy.
// It retrieves the toy ID from the request parameters, parses it,
// and then calls the ToysRepository to get the toy by its ID.
// If an error occurs during the retrieval, it returns a 500 HTTP error.
// Finally, it renders the redeemed toy using the lib.Render function.
func (h RedeemToyHandler) RedeemToyViewHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	toy, err := h.ToysRepository.GetToyByID(toyID)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return lib.Render(ctx, redeem.RedeemToy(toy))
}
