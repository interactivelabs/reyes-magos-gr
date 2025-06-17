package admin

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	toys_view "reyes-magos-gr/views/admin/toys"
	"strconv"
	"strings"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type ToysHandler struct {
	ToysRepository store.ToysRepository
}

func (h ToysHandler) ToysViewHandler(ctx echo.Context) error {
	toys, err := h.ToysRepository.GetToys()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, toys_view.Toys(toys))
}

func (h ToysHandler) CreateToyFormHandler(ctx echo.Context) error {
	return lib.Render(ctx, toys_view.CreateToyForm())
}

type CreateToyRequest struct {
	ToyName        string `form:"toy_name" validate:"required"`
	ToyDescription string `form:"toy_description"`
	Category       string `form:"category" validate:"required"`
	AgeMin         int64  `form:"age_min" validate:"required,min=1,max=16"`
	AgeMax         int64  `form:"age_max" validate:"required,max=16"`
	Image1         string `form:"image1" validate:"required,url"`
	Image2         string `form:"image2" validate:"required,url"`
	Image3         string `form:"image3" validate:"required,url"`
	SourceURL      string `form:"source_url" validate:"required,url"`
}

func (h ToysHandler) CreateToyPostHandler(ctx echo.Context) error {
	tr := new(CreateToyRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var toy models.Toy
	err := dto.Map(&toy, tr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	toyID, err := h.ToysRepository.CreateToy(toy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	toy.ToyID = toyID

	return lib.Render(ctx, toys_view.ToyRow(toy))
}

func (h ToysHandler) UpdateToyFormHandler(ctx echo.Context) error {
	toyID, err := strconv.ParseInt(ctx.Param("toy_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	toy, err := h.ToysRepository.GetToyByID(toyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, toys_view.UpdateToyForm(toy))
}

func (h ToysHandler) UpdateToyPutHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tr := new(CreateToyRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var toy models.Toy
	err = dto.Map(&toy, tr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	toy.ToyID = toyID

	var categories = strings.Split(toy.Category, ",")
	if len(categories) > 0 {
		for i, category := range categories {
			categories[i] = strings.TrimSpace(category)
		}
		toy.Category = strings.Join(categories, ",")
	}

	err = h.ToysRepository.UpdateToy(toy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, toys_view.ToyRow(toy))
}

func (h ToysHandler) DeleteToyHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.ToysRepository.DeleteToy(toyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.NoContent(http.StatusOK)
}

type SearchBoxItem struct {
	Value string
	Label string
}

func (h ToysHandler) ToysCategoriesViewHandler(ctx echo.Context) error {
	categories, err := h.ToysRepository.GetCategories()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var objects []SearchBoxItem
	for _, str := range categories {
		obj := SearchBoxItem{Value: str, Label: str}
		objects = append(objects, obj)
	}

	return ctx.JSON(http.StatusOK, objects)
}
