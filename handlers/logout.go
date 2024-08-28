package handlers

import (
	"net/http"
	"net/url"
	"os"
	"reyes-magos-gr/lib"

	"github.com/labstack/echo/v4"
)

func (h LoginHandler) LogoutRedirectHandler(ctx echo.Context) error {
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
