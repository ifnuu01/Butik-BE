package utils

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			errors := make(map[string]string)
			val := reflect.ValueOf(i)

			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			typ := val.Type()
			for _, ve := range validationErrors {
				field, _ := typ.FieldByName(ve.Field())
				jsonTag := field.Tag.Get("json")
				jsonName := strings.Split(jsonTag, ",")[0]
				if jsonName == "" {
					jsonName = strings.ToLower(ve.Field())
				}
				errors[jsonName] = ve.Tag()
			}
			return echo.NewHTTPError(http.StatusBadRequest, map[string]interface{}{
				"status":  400,
				"message": "validation error",
				"errors":  errors,
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}
	return nil
}

func ValidaionErrorResponse(c echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		return c.JSON(httpErr.Code, httpErr.Message)
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
