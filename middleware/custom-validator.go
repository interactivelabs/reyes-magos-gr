package middleware

import (
	"reflect"
	"reyes-magos-gr/lib"
	"time"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
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

	return err == nil
}

func registerCustomFieldRules(v *validator.Validate) {
	v.RegisterValidation("iso_8601_date", iso8601Date)
}

func registerCustomTags(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if tag := fld.Tag.Get("json"); tag != "" {
			return tag
		}
		if tag := fld.Tag.Get("form"); tag != "" {
			return tag
		}
		return fld.Name
	})
}

func NewValidator() *CustomValidator {
	v := validator.New()
	registerCustomFieldRules(v)
	registerCustomTags(v)
	return &CustomValidator{validator: v}
}
