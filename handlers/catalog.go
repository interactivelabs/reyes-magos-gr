package handlers

import (
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	catalog "reyes-magos-gr/views/catalog"

	"github.com/labstack/echo/v4"
)

type CatalogHandler struct {
	ToysRepository repository.ToysRepository
}

func (h CatalogHandler) CatalogViewHandler(ctx echo.Context) error {
	toys, err := h.ToysRepository.GetToys()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return lib.Render(ctx, catalog.Catalog(toys))
}
