package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/platform/authenticator"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	Auth *authenticator.Authenticator
}

func NewLoginHandler(auth *authenticator.Authenticator) *LoginHandler {
	return &LoginHandler{
		Auth: auth,
	}
}

func (h *LoginHandler) LoginRedirectHandler(ctx echo.Context) error {
	state, err := generateRandomState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := lib.SetCookieSession(ctx, "state_session", state); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, h.Auth.AuthCodeURL(state))
}

func (h *LoginHandler) LoginCallbackHandler(ctx echo.Context) error {
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

	var profile map[string]any
	if err := idToken.Claims(&profile); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if err := lib.SetSessionProfile(ctx, profile); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	stateSession.Options.MaxAge = -1
	if err := stateSession.Save(ctx.Request(), ctx.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Redirect to logged in page.
	return ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func (h *LoginHandler) LogoutRedirectHandler(ctx echo.Context) error {
	if err := lib.DeleteCookieSession(ctx, "access_token"); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := lib.DeleteCookieSession(ctx, "profile"); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	scheme := "https"
	env := os.Getenv("ENV")
	if env == "development" {
		scheme = "http"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request().Host)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	return ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
