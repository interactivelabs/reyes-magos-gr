package handlers

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	catalog "reyes-magos-gr/views/catalog"

	"github.com/labstack/echo/v4"
)

type CatalogHandler struct {
	ToysRepository repository.ToysRepository
}

type CatalogHandlerRequest struct {
	AgeMin   int64    `query:"age_min"`
	AgeMax   int64    `query:"age_max"`
	Category []string `query:"category"`
}

func (h CatalogHandler) CatalogViewHandler(ctx echo.Context) error {
	cr := new(CatalogHandlerRequest)
	if err := ctx.Bind(cr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	toys, err := h.ToysRepository.GetToysWithFilters(cr.AgeMin, cr.AgeMax, cr.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	categories, err := h.ToysRepository.GetCategories()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, catalog.Catalog(toys, categories, cr.Category, int(cr.AgeMin), int(cr.AgeMax)))
}
