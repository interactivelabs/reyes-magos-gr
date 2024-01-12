package handlers

import (
	"database/sql"

	redeem "reyes-magos-gr/components/redeem"

	"github.com/labstack/echo/v4"
)

type RedeemHandler struct {
	DB *sql.DB
}

func (h RedeemHandler) RedeemViewHandler(ctx echo.Context) error {
	return render(ctx, redeem.Redeem())
}
