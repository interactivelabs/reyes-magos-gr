package handlers

import (
	"reyes-magos-gr/lib"
	redeem "reyes-magos-gr/views/redeem"

	"github.com/labstack/echo/v4"
)

type RedeemHandler struct {
}

func (h RedeemHandler) RedeemViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, redeem.Redeem())
}
