package admin

import (
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	toys_view "reyes-magos-gr/views/admin/toys"
	"strconv"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type ToysHandler struct {
	ToysRepository repository.ToysRepository
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
		return err
	}

	var toy model.Toy
	err := dto.Map(&toy, tr)
	if err != nil {
		return err
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
		return err
	}

	var toy model.Toy
	err = dto.Map(&toy, tr)
	if err != nil {
		return err
	}

	toy.ToyID = toyID

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
