package api

import (
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type CreaetToyRequest struct {
	ToyName   string `json:"toy_name" validate:"required"`
	AgeMin    int64  `json:"age_min" validate:"required min=1,max=16"`
	AgeMax    int64  `json:"age_max" validate:"required max=16"`
	Image1    string `json:"image1" validate:"required url"`
	Image2    string `json:"image2" validate:"required url"`
	SourceURL string `json:"source_url" validate:"required url"`
}

type ToyHandler struct {
	ToysRepository repository.ToysRepository
}

func (h ToyHandler) CreateToyApiHandler(ctx echo.Context) error {
	tr := new(CreaetToyRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return err
	}

	var toy model.Toy
	dto.Map(&toy, tr)

	toyID, err := h.ToysRepository.CreateToy(toy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, toyID)
}

type UpdateToyRequest struct {
	ToyID     int64  `json:"toy_id" validate:"required"`
	ToyName   string `json:"toy_name"`
	AgeMin    int64  `json:"age_min" validate:"omitempty,min=1,max=16"`
	AgeMax    int64  `json:"age_max" validate:"omitempty,max=16"`
	Image1    string `json:"image1" validate:"omitempty,url"`
	Image2    string `json:"image2" validate:"omitempty,url"`
	SourceURL string `json:"source_url" validate:"omitempty,url"`
}

func (h ToyHandler) UpdateToyApiHandler(ctx echo.Context) error {
	tr := new(UpdateToyRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(tr); err != nil {
		return err
	}

	var toy model.Toy
	dto.Map(&toy, tr)

	err := h.ToysRepository.UpdateToy(toy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, toy.ToyID)
}
