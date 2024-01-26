package handlers

import (
	"reyes-magos-gr/lib"
	home "reyes-magos-gr/views/home"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) HomeViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, home.Home())
}
