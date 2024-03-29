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
	isAdmin := IsAdmin(ctx)
	c := context.WithValue(ctx.Request().Context(), isAdminKey, isAdmin)
	return component.Render(c, ctx.Response())
}

func IsAdmin(ctx echo.Context) bool {
	user := ctx.Get("user")
	if user == nil {
		return false
	}

	token := user.(*jwt.Token)
	if token.Claims == nil {
		return false
	}

	if claims, ok := token.Claims.(*api.JwtCustomClaims); ok {
		return claims.Admin
	}

	return false
}

func GetIsAdmin(ctx context.Context) bool {
	if isAdmin, ok := ctx.Value(isAdminKey).(bool); ok {
		return isAdmin
	}
	return false
}
