package api

import (
	"math/rand"
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type CodeHandler struct {
	CodesRepository repository.CodesRepository
}

type CreateCodeResult struct {
	Code   string `json:"code"`
	CodeID int64  `json:"code_id"`
}

func (h CodeHandler) CreateCodeApiHandler(ctx echo.Context) error {
	code := model.Code{
		Code:       generateRandomString(6),
		Expiration: time.Now().AddDate(0, 0, 10).Format("2006-01-02"),
	}

	_, codeRow, err := h.CodesRepository.CreateCode(code)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	result := CreateCodeResult{
		Code:   codeRow.Code,
		CodeID: codeRow.CodeID,
	}

	return ctx.JSON(http.StatusOK, result)
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

	var codes []CreateCodeResult
	for i := 0; i < int(cr.Count); i++ {
		code := model.Code{
			Code:       generateRandomString(6),
			Expiration: time.Now().AddDate(0, 0, 10).Format("2006-01-02"),
		}

		_, codeRow, err := h.CodesRepository.CreateCode(code)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		result := CreateCodeResult{
			Code:   codeRow.Code,
			CodeID: codeRow.CodeID,
		}

		codes = append(codes, result)
	}

	result := CreateCodeBatchResult{
		Codes: codes,
	}

	return ctx.JSON(http.StatusOK, result)
}

func generateRandomString(length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(chars[rand.Intn(len(chars))])
	}
	return sb.String()
}
