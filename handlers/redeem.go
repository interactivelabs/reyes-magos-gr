package handlers

import (
	"github.com/labstack/echo/v4"
	"reyes-magos-gr/components/redeem"
)

type RedeemHandler struct {}

func (h RedeemHandler) RedeemViewHandler(ctx echo.Context) error {
	return render(ctx, redeem.Main())
}