package middleware

import (
	"encoding/base64"
	"encoding/json"
	"maps"
	"net/http"
	"reflect"
	"reyes-magos-gr/lib"

	"github.com/labstack/echo/v4"
	"github.com/posthog/posthog-go"
)

var Flags = [1]string{"volunteers-cart-enabled"}

func UserFlags(client posthog.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			sessionProfile, err := lib.GetSessionProfile(ctx)
			if err != nil || sessionProfile == nil {
				return next(ctx)
			}

			email, _ := sessionProfile["name"].(string)
			DistinctId := base64.StdEncoding.EncodeToString([]byte(email))

			profileFlags := make(map[string]bool)
			if sessionFlags, ok := sessionProfile["flags"].(string); ok {
				if err := json.Unmarshal([]byte(sessionFlags), &profileFlags); err != nil {
					profileFlags = make(map[string]bool)
				}
			}

			originalFlags := make(map[string]bool)
			maps.Copy(originalFlags, profileFlags)

			for _, flag := range Flags {
				// If the flag is already set in the session, skip checking it again
				if _, ok := profileFlags[flag]; ok {
					continue
				}

				isFlagEnabled, err := client.IsFeatureEnabled(
					posthog.FeatureFlagPayload{
						Key:        flag,
						DistinctId: DistinctId,
					})

				if err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}

				profileFlags[flag] = isFlagEnabled == true
			}

			if reflect.DeepEqual(originalFlags, profileFlags) == false {
				if flags, err := json.Marshal(profileFlags); err == nil {
					sessionProfile["flags"] = string(flags)
				}
				if err := lib.SetSessionProfile(ctx, sessionProfile); err != nil {
					return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
				}
			}

			if err := next(ctx); err != nil {
				return err
			}

			return nil
		}
	}
}
