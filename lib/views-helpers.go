package lib

import (
	"context"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, component templ.Component) error {
	profileView := GetProfileView(ctx)
	c := context.WithValue(ctx.Request().Context(), profileKey, profileView)
	return component.Render(c, ctx.Response())
}

func FormatDate(date string) (string, error) {
	fmtDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", err
	}
	return fmtDate.Format("January 2, 2006"), nil
}
