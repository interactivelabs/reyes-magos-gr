package handlers

import (
	redeem "reyes-magos-gr/views/redeem"

	"github.com/labstack/echo/v4"
)

type RedeemHandler struct {
}

func (h RedeemHandler) RedeemViewHandler(ctx echo.Context) error {
	return render(ctx, redeem.Redeem())
}
