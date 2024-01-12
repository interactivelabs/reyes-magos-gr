package handlers

import (
	codes "reyes-magos-gr/components/codes"

	"github.com/labstack/echo/v4"
)

type CodesHandler struct{}

func (h CodesHandler) CodesViewHandler(ctx echo.Context) error {
	return render(ctx, codes.Codes())
}
