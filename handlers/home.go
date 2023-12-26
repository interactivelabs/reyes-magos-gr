package handlers

import (
	home "reyes-magos-gr/components/home"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) HomeViewHandler(ctx echo.Context) error {
	return render(ctx, home.Home())
}
