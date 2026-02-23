package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DriverStatus string

const (
	StatusAvailable DriverStatus = "AVAILABLE"
	StatusOnTrip    DriverStatus = "ON_TRIP"
	StatusOffline   DriverStatus = "OFFLINE"
)

type Driver struct {
	ID        uuid.UUID    `json:"id"`
	Name      string       `json:"name"`
	Status    DriverStatus `json:"status"`
	CreatedAt time.Time    `json:"createdAt"`
	UpdatedAt time.Time    `json:"updatedAt"`
}

type DriverRepository interface {
	Create(ctx context.Context, driver *Driver) error
	GetByID(ctx context.Context, id uuid.UUID) (*Driver, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status DriverStatus) error
}
