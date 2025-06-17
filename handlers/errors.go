package handlers

import (
	"reyes-magos-gr/lib"
	errors "reyes-magos-gr/views/errors"

	"github.com/labstack/echo/v4"
)

func (h *HomeHandler) Error401(ctx echo.Context) error {
	return lib.Render(ctx, errors.Error401())
}

func (h *HomeHandler) Error404(ctx echo.Context) error {
	return lib.Render(ctx, errors.Error404())
}

func (h *HomeHandler) Error500(ctx echo.Context) error {
	return lib.Render(ctx, errors.Error500())
}
