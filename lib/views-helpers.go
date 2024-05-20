package lib

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, component templ.Component) error {
	profileView := GetProfileView(ctx)
	c := context.WithValue(ctx.Request().Context(), profileKey, profileView)
	return component.Render(c, ctx.Response())
}
