package middleware

import (
	"encoding/json"
	"net/http"
	"reyes-magos-gr/lib"

	"github.com/labstack/echo/v4"
	"github.com/posthog/posthog-go"
)

func UserFlags(client posthog.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			profile := lib.GetCtxProfile(ctx)

			if profile.Email == "" {
				return next(ctx)
			}

			isMyFlagEnabled, err := client.IsFeatureEnabled(
				posthog.FeatureFlagPayload{
					Key:        "volunteers-cart-enabled",
					DistinctId: "123", // Replace with actual user ID or distinct ID
				})

			if isMyFlagEnabled == true {
				// Do something differently for this user
			}

			profileJSON, err := json.Marshal(profile)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			if err := lib.SetCookieSession(ctx, "profile", string(profileJSON)); err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			if err := next(ctx); err != nil {
				return err
			}

			return next(ctx)
		}
	}
}
