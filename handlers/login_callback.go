package handlers

import (
	"encoding/json"
	"net/http"
	"reyes-magos-gr/lib"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func (h LoginHandler) LoginCallbackHandler(ctx echo.Context) error {
	stateSession, err := session.Get("state_session", ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	state := stateSession.Values["state_session"]
	if state != ctx.QueryParam("state") {
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to authenticate")
	}

	code := ctx.QueryParam("code")
	token, err := h.Auth.Exchange(ctx.Request().Context(), code)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to authenticate")
	}

	if err := lib.SetCookieSession(ctx, "access_token", token.AccessToken); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	idToken, err := h.Auth.VerifyIDToken(ctx.Request().Context(), token)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Failed to authenticate")
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := lib.SetCookieSession(ctx, "profile", string(profileJSON)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	stateSession.Options.MaxAge = -1
	if err := stateSession.Save(ctx.Request(), ctx.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to logged in page.
	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}
