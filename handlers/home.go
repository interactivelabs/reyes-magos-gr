package handlers

import (
	"reyes-magos-gr/lib"
	home "reyes-magos-gr/views/home"
	support "reyes-magos-gr/views/support"

	"github.com/labstack/echo/v4"
)

type HomeHandler struct{}

func (h HomeHandler) HomeViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, home.Home())
}

func (h HomeHandler) SupportViewHandler(ctx echo.Context) error {
	return lib.Render(ctx, support.Support())
}
