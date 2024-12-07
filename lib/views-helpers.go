package lib

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(ctx echo.Context, component templ.Component) error {
	profileView := GetProfileView(ctx)
	c := ctx.Request().Context()
	c = context.WithValue(c, profileKey, profileView)
	return component.Render(c, ctx.Response())
}

func FormatDate(date string) (string, error) {
	fmtDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}
	return fmtDate.Format(TextDate), nil
}

func GetSafeIdUrl(url string, id int64) string {
	return string(templ.URL(fmt.Sprintf(url, id)))
}

func GetAssetUrl(url string) string {
	env := os.Getenv("ENV")
	if env == "development" {
		return string(templ.URL(fmt.Sprintf("/public/%s", url)))
	}
	return string(templ.URL(fmt.Sprintf("https://static.dl-toys.com/%s", url)))
}

func GetHTMLErrorCode(err error) (code int) {
	if he, ok := err.(*echo.HTTPError); ok {
		return he.Code
	}
	return 0
}

func GetPaginationLink(currentQuery string, page int, pageSize int64) string {
	if currentQuery == "" {
		return fmt.Sprintf("/catalog?page=%d&page_size=%d", page, pageSize)
	}
	return fmt.Sprintf("/catalog?page=%d&page_size=%d&%s", page, pageSize, currentQuery)
}

func ReadSylesFromFile() string {
	styles, err := os.ReadFile("/public/css/main.css")
	if err != nil {
		// handle error
		return ""
	}

	return string(styles)
}
