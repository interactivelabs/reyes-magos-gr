package middleware

import (
	"net/http"

	"reyes-magos-gr/lib"
	errors "reyes-magos-gr/views/errors"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, c echo.Context) {

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	c.Logger().Error(err)

	var errorPage func() templ.Component

	switch code {
	case 401:
		errorPage = errors.Error401
	case 404:
		errorPage = errors.Error404
	case 500:
		errorPage = errors.Error500
	}

	lib.Render(c, errorPage())

}
