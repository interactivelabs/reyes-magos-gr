package lib

import (
	"context"
	"reyes-magos-gr/api"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type contextKey string

const (
	isAdminKey contextKey = "isAdmin"
)

func Render(ctx echo.Context, component templ.Component) error {
	user := ctx.Get("user").(*jwt.Token)
	isAdmin := false
	if claims, ok := user.Claims.(*api.JwtCustomClaims); ok {
		isAdmin = claims.Admin
	}

	c := context.WithValue(ctx.Request().Context(), isAdminKey, isAdmin)

	return component.Render(c, ctx.Response())
}

func GetIsAdmin(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(isAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}
