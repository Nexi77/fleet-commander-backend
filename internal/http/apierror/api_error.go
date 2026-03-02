package apierror

import "github.com/Nexi77/fleetcommander-backend/internal/http/validation"

type Code string

const (
	CodeBadRequestInvalidJSON Code = "BAD_REQUEST_INVALID_JSON"
	CodeBadRequestValidation  Code = "BAD_REQUEST_VALIDATION"
	CodeDriverCreateFailed    Code = "DRIVER_CREATE_FAILED"
)

type Response struct {
	Code    Code                              `json:"code"`
	Details []validation.FieldValidationError `json:"details,omitempty"`
}
