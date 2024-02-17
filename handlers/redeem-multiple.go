package handlers

import (
	"reyes-magos-gr/lib"

	redeemMultiple "reyes-magos-gr/views/redeem-multiple"

	"github.com/labstack/echo/v4"
)

type RedeemMultipleHandler struct {
}

func (h RedeemMultipleHandler) RedeemMultipleViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, redeemMultiple.RedeemMultiple())
}
