package api

import (
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"strconv"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type CreateToyRequest struct {
	ToyName        string `json:"toy_name" validate:"required"`
	ToyDescription string `json:"toy_description" validate:"required"`
	Category       string `json:"category" validate:"required"`
	AgeMin         int64  `json:"age_min" validate:"required,min=1,max=16"`
	AgeMax         int64  `json:"age_max" validate:"required,max=16"`
	Image1         string `json:"image1" validate:"required,url"`
	Image2         string `json:"image2" validate:"omitempty,url"`
	Image3         string `json:"image3" validate:"omitempty,url"`
	SourceURL      string `json:"source_url" validate:"required,url"`
}

type ToyHandler struct {
	ToysRepository repository.ToysRepository
}

func (h ToyHandler) CreateToyApiHandler(ctx echo.Context) error {
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

	return ctx.JSON(http.StatusOK, toyID)
}

func (h ToyHandler) CreateBatchToysApiHandler(ctx echo.Context) error {
	tr := new([]CreateToyRequest)
	if err := ctx.Bind(tr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var toys []model.Toy
	err := dto.Map(&toys, tr)
	if err != nil {
		return err
	}

	var toysIDs []int64
	for _, toy := range toys {
		toyID, err := h.ToysRepository.CreateToy(toy)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		toysIDs = append(toysIDs, toyID)
	}

	return ctx.JSON(http.StatusOK, toysIDs)
}

type UpdateToyRequest struct {
	ToyID          int64  `json:"toy_id" validate:"required"`
	ToyName        string `json:"toy_name"`
	ToyDescription string `json:"toy_description"`
	Category       string `json:"category"`
	AgeMin         int64  `json:"age_min" validate:"omitempty,min=1,max=16"`
	AgeMax         int64  `json:"age_max" validate:"omitempty,max=16"`
	Image1         string `json:"image1" validate:"omitempty,url"`
	Image2         string `json:"image2" validate:"omitempty,url"`
	Image3         string `json:"image3" validate:"omitempty,url"`
	SourceURL      string `json:"source_url" validate:"omitempty,url"`
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
	err := dto.Map(&toy, tr)
	if err != nil {
		return err
	}

	err = h.ToysRepository.UpdateToy(toy)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, toy.ToyID)
}

type DeleteToyRequest struct {
}

func (h ToyHandler) DeleteToyApiHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid toy ID")
	}

	err = h.ToysRepository.DeleteToy(toyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, toyID)
}
