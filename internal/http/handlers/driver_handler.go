package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/Nexi77/fleetcommander-backend/internal/domain"
	"github.com/Nexi77/fleetcommander-backend/internal/http/apierror"
	"github.com/Nexi77/fleetcommander-backend/internal/http/response"
	"github.com/Nexi77/fleetcommander-backend/internal/http/validation"
	playgroundValidator "github.com/go-playground/validator/v10"
)

const (
	maxRequestBodyBytes int64 = 1 << 20
)

type DriverHandler struct {
	driverRepo domain.DriverRepository
	validate   *playgroundValidator.Validate
}

func NewDriverHandler(driverRepo domain.DriverRepository, validate *playgroundValidator.Validate) *DriverHandler {
	if validate == nil {
		validate = validation.NewValidator()
	}

	return &DriverHandler{
		driverRepo: driverRepo,
		validate:   validate,
	}
}

type registerDriverRequest struct {
	Name   string `json:"name" validate:"required,min=2,max=100"`
	Status string `json:"status" validate:"omitempty,oneof=AVAILABLE ON_TRIP OFFLINE"`
}

type driverResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (h *DriverHandler) RegisterDriver(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req registerDriverRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, apierror.Response{Code: apierror.CodeBadRequestInvalidJSON})
		return
	}

	if err := decoder.Decode(&struct{}{}); err != nil && err != io.EOF {
		response.WriteJSON(w, http.StatusBadRequest, apierror.Response{Code: apierror.CodeBadRequestInvalidJSON})
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Status = strings.ToUpper(strings.TrimSpace(req.Status))

	validationErrors := validation.ValidateDTO(h.validate, req)
	if len(validationErrors) > 0 {
		response.WriteJSON(w, http.StatusBadRequest, apierror.Response{
			Code:    apierror.CodeBadRequestValidation,
			Details: validationErrors,
		})
		return
	}

	driver := &domain.Driver{
		Name:   req.Name,
		Status: parseDriverStatus(req.Status),
	}

	if err := h.driverRepo.Create(ctx, driver); err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, apierror.Response{Code: apierror.CodeDriverCreateFailed})
		return
	}

	response.WriteJSON(w, http.StatusCreated, driverResponse{
		ID:        driver.ID.String(),
		Name:      driver.Name,
		Status:    string(driver.Status),
		CreatedAt: driver.CreatedAt.UTC().Format(http.TimeFormat),
		UpdatedAt: driver.UpdatedAt.UTC().Format(http.TimeFormat),
	})
}

func parseDriverStatus(raw string) domain.DriverStatus {
	if strings.TrimSpace(raw) == "" {
		return domain.StatusOffline
	}

	status := domain.DriverStatus(raw)
	switch status {
	case domain.StatusAvailable, domain.StatusOnTrip, domain.StatusOffline:
		return status
	default:
		return domain.StatusOffline
	}
}
