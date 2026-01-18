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
	// Sanitize
	SanitizeStruct(i)

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
				errors[jsonName] = formatValidationError(ve)
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

// SanitizeStruct trims whitespace
func SanitizeStruct(i interface{}) {
	val := reflect.ValueOf(i)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for j := 0; j < val.NumField(); j++ {
		field := val.Field(j)
		if field.Kind() == reflect.String && field.CanSet() {
			trimmed := strings.TrimSpace(field.String())
			field.SetString(trimmed)
		}
	}
}

// formatValidationError
func formatValidationError(ve validator.FieldError) string {
	switch ve.Tag() {
	case "required":
		return "field is required"
	case "min":
		return "minimum length is " + ve.Param()
	case "max":
		return "maximum length is " + ve.Param()
	case "gt":
		return "must be greater than " + ve.Param()
	case "gte":
		return "must be greater than or equal to " + ve.Param()
	case "lt":
		return "must be less than " + ve.Param()
	case "lte":
		return "must be less than or equal to " + ve.Param()
	case "email":
		return "must be a valid email"
	case "numeric":
		return "must be numeric"
	case "alphanum":
		return "must be alphanumeric"
	case "oneof":
		return "must be one of: " + ve.Param()
	default:
		return ve.Tag()
	}
}

func ValidationErrorResponse(c echo.Context, err error) error {
	if httpErr, ok := err.(*echo.HTTPError); ok {
		return c.JSON(httpErr.Code, httpErr.Message)
	}
	return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
}
