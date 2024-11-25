package handlers

import (
	"fmt"
	"math"
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/handlers/dtos"
	"reyes-magos-gr/lib"
	catalog "reyes-magos-gr/views/catalog"

	"github.com/labstack/echo/v4"
)

const PageSize int64 = 12

type CatalogHandler struct {
	ToysRepository repository.ToysRepository
}

type CatalogHandlerRequest struct {
	AgeMin   int64    `query:"age_min"`
	AgeMax   int64    `query:"age_max"`
	Category []string `query:"category"`
	Page     int64    `query:"page"`
	PageSize int64    `query:"page_size"`
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

	pageSize := PageSize
	if cr.PageSize > 0 {
		pageSize = cr.PageSize
	}

	toys, err := h.ToysRepository.GetToysWithFiltersPaged(page, pageSize, cr.AgeMin, cr.AgeMax, cr.Category)
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

	pagesFloat := float64(count) / float64(pageSize)
	pages := int64(math.Ceil(pagesFloat))

	currentQueryVlues := ctx.Request().URL.Query()
	currentQueryVlues.Del("page")
	currentQueryVlues.Del("page_size")
	currentQuery := currentQueryVlues.Encode()

	links := make([]dtos.CatalogLink, 0)
	for i := 1; i <= int(pages); i++ {
		links = append(links, dtos.CatalogLink{
			Text:  fmt.Sprintf("%d", i),
			Url:   lib.GetPaginationLink(currentQuery, i, pageSize),
			Label: fmt.Sprintf("Ir a la pagina %d de %d", i, pages),
		})
	}

	return lib.Render(
		ctx,
		catalog.Catalog(
			toys,
			categories,
			page,
			pages,
			pageSize,
			count,
			currentQuery,
			cr.Category,
			int64(cr.AgeMin),
			int64(cr.AgeMax),
			links,
		))
}
