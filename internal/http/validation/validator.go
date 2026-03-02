package validation

import (
	"errors"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type ValidationCode string

const (
	CodeValidationInvalidSchema ValidationCode = "VALIDATION_INVALID_SCHEMA"
	CodeValidationFailed        ValidationCode = "VALIDATION_FAILED"
	CodeValidationRequired      ValidationCode = "VALIDATION_REQUIRED"
	CodeValidationOneOf         ValidationCode = "VALIDATION_ONE_OF"
	CodeValidationMin           ValidationCode = "VALIDATION_MIN"
	CodeValidationMax           ValidationCode = "VALIDATION_MAX"
	CodeValidationInvalid       ValidationCode = "VALIDATION_INVALID"
)

const (
	fieldRequest = "request"
	tagRequired  = "required"
	tagOneOf     = "oneof"
	tagMin       = "min"
	tagMax       = "max"
)

type FieldValidationError struct {
	Field  string            `json:"field"`
	Code   ValidationCode    `json:"code"`
	Params map[string]string `json:"params,omitempty"`
}

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tag := field.Tag.Get("json")
		if tag == "" {
			return field.Name
		}

		name := strings.Split(tag, ",")[0]
		if name == "" || name == "-" {
			return field.Name
		}

		return name
	})

	return v
}

func ValidateDTO(v *validator.Validate, dto any) []FieldValidationError {
	if err := v.Struct(dto); err != nil {
		var invalidValidationErr *validator.InvalidValidationError
		if errors.As(err, &invalidValidationErr) {
			return []FieldValidationError{
				{
					Field: fieldRequest,
					Code:  CodeValidationInvalidSchema,
				},
			}
		}

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			out := make([]FieldValidationError, 0, len(validationErrs))

			for _, fieldErr := range validationErrs {
				out = append(out, FieldValidationError{
					Field:  fieldErr.Field(),
					Code:   mapValidationCode(fieldErr.Tag()),
					Params: buildErrorParams(fieldErr),
				})
			}

			return out
		}

		return []FieldValidationError{
			{
				Field: fieldRequest,
				Code:  CodeValidationFailed,
			},
		}
	}

	return nil
}

func mapValidationCode(tag string) ValidationCode {
	switch tag {
	case tagRequired:
		return CodeValidationRequired
	case tagOneOf:
		return CodeValidationOneOf
	case tagMin:
		return CodeValidationMin
	case tagMax:
		return CodeValidationMax
	default:
		return CodeValidationInvalid
	}
}

func buildErrorParams(fieldErr validator.FieldError) map[string]string {
	if fieldErr.Param() == "" {
		return nil
	}

	return map[string]string{
		"param": fieldErr.Param(),
	}
}
