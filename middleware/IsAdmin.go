package middleware

import (
	"reyes-magos-gr/lib"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func IsAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {

			accessTokenSession, err := session.Get("access_token", ctx)
			if err != nil {
				return echo.ErrUnauthorized
			}

			token := accessTokenSession.Values["access_token"]
			if token == nil {
				return echo.ErrUnauthorized
			}

			profile, err := lib.GetProfile(ctx)
			if err != nil || profile == nil {
				return echo.ErrUnauthorized
			}

			return next(ctx)
		}
	}
}
