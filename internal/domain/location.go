package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// DriverLocation represents a single GPS ping from a driver's mobile app.
type DriverLocation struct {
	DriverID  uuid.UUID `json:"driverId"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// LocationRepository defines the contract for high-throughput, real-time location data (Hot Data).
type LocationRepository interface {
	SaveLocation(ctx context.Context, location *DriverLocation) error
	GetNearbyDrivers(ctx context.Context, lat, lon float64, radiusKm float64) ([]DriverLocation, error)
}
