package handlers

import (
	"math"
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	catalog "reyes-magos-gr/views/catalog"

	"github.com/labstack/echo/v4"
)

const PAGE_SIZE int64 = 12

type CatalogHandler struct {
	ToysRepository repository.ToysRepository
}

type CatalogHandlerRequest struct {
	AgeMin   int64    `query:"age_min"`
	AgeMax   int64    `query:"age_max"`
	Category []string `query:"category"`
	Page     int64    `query:"page"`
}

func (h CatalogHandler) CatalogViewHandler(ctx echo.Context) error {
	cr := new(CatalogHandlerRequest)
	if err := ctx.Bind(cr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	page := cr.Page
	if page < 1 {
		page = 1
	}

	toys, err := h.ToysRepository.GetToysWithFiltersPaged(page, PAGE_SIZE, cr.AgeMin, cr.AgeMax, cr.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	count, err := h.ToysRepository.GetToysCountWithFilters(cr.AgeMin, cr.AgeMax, cr.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	categories, err := h.ToysRepository.GetCategories()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	pagesFloat := float64(count) / float64(PAGE_SIZE)
	pages := int64(math.Ceil(pagesFloat))

	currentQueryVlues := ctx.Request().URL.Query()
	currentQueryVlues.Del("page")
	currentQuery := currentQueryVlues.Encode()

	return lib.Render(
		ctx,
		catalog.Catalog(toys, categories, page, pages, PAGE_SIZE, count, currentQuery, cr.Category, int64(cr.AgeMin), int64(cr.AgeMax)))
}
