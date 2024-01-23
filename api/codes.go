package api

import (
	"net/http"
	"reyes-magos-gr/services"

	"github.com/dranikpg/dto-mapper"
	"github.com/labstack/echo/v4"
)

type CodeHandler struct {
	CodesService services.CodesService
}

type CreateCodeResult struct {
	Code   string `json:"code"`
	CodeID int64  `json:"code_id"`
}

func (h CodeHandler) CreateCodeApiHandler(ctx echo.Context) error {
	code, err := h.CodesService.CreateCode()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var codesResult CreateCodeResult
	dto.Map(&codesResult, code)

	return ctx.JSON(http.StatusOK, codesResult)
}

// new handler to create several codes at once
type CreateCodeBatchRequest struct {
	Count int64 `json:"count" validate:"required,min=1,max=100"`
}

type CreateCodeBatchResult struct {
	Codes []CreateCodeResult `json:"codes"`
}

func (h CodeHandler) CreateCodeBatchApiHandler(ctx echo.Context) error {
	cr := new(CreateCodeBatchRequest)
	if err := ctx.Bind(cr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(cr); err != nil {
		return err
	}

	codes, err := h.CodesService.CreateCodeBatch(cr.Count)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var codesResult []CreateCodeResult
	dto.Map(&codesResult, codes)

	result := CreateCodeBatchResult{
		Codes: codesResult,
	}

	return ctx.JSON(http.StatusOK, result)
}
