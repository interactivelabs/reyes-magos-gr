package middleware

import (
	"net/http"
	"reyes-magos-gr/lib"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func iso8601Date(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	// only validate if the field is not empty, required is a different rule
	if value == "" {
		return true
	}

	_, err := time.Parse(lib.YYYYMMDD, fl.Field().String())
	if err != nil {
		return false
	}

	return true
}

func registerCustomFieldRules(v *validator.Validate) {
	v.RegisterValidation("iso_8601_date", iso8601Date)
}

func NewValidator() *CustomValidator {
	v := validator.New()
	registerCustomFieldRules(v)
	return &CustomValidator{validator: v}
}
