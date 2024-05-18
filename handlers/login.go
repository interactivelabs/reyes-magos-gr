package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/platform/authenticator"

	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	Auth *authenticator.Authenticator
}

func (h LoginHandler) LoginRedirectHandler(ctx echo.Context) error {
	state, err := generateRandomState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if err := lib.SetCookieSession(ctx, "state_session", state); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, h.Auth.AuthCodeURL(state))
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
