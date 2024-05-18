package lib

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	isAdminKey contextKey = "isAdmin"
)

func Render(ctx echo.Context, component templ.Component) error {
	isAdmin := IsAdmin(ctx)
	c := context.WithValue(ctx.Request().Context(), isAdminKey, isAdmin)
	return component.Render(c, ctx.Response())
}

func IsAdmin(ctx echo.Context) bool {
	user := ctx.Get("user")
	if user == nil {
		return false
	}

	return false
}

func GetIsAdmin(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(isAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}
