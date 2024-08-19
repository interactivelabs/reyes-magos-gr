package admin

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	toys_view "reyes-magos-gr/views/admin/toys"

	"github.com/labstack/echo/v4"
)

type ToysHandler struct {
	ToysRepository repository.ToysRepository
}

func (h ToysHandler) ToysViewHandler(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	toys, err := h.ToysRepository.GetToys()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, toys_view.Toys(toys))
}

func (h ToysHandler) CreateToyFormHandler(ctx echo.Context) error {
	profile := lib.GetProfileView(ctx)
	if profile.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	return lib.Render(ctx, toys_view.CreateToyForm())
}
