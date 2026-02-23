package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/Nexi77/fleetcommander-backend/internal/domain"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	locationsKey  = "drivers:locations"
	timestampsKey = "drivers:timestamps"
)

type locationRepository struct {
	client *redis.Client
}

func NewLocationRepository(client *redis.Client) domain.LocationRepository {
	return &locationRepository{
		client: client,
	}
}

func (r *locationRepository) SaveLocation(ctx context.Context, loc *domain.DriverLocation) error {
	err := r.client.GeoAdd(ctx, locationsKey, &redis.GeoLocation{
		Name:      loc.DriverID.String(),
		Longitude: loc.Longitude,
		Latitude:  loc.Latitude,
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to save driver location in redis: %w", err)
	}

	// Optional but recommended for Hot Data:
	// Store the exact timestamp in a separate Hash to know EXACTLY when they last pinged.
	// HSET drivers:timestamps <driver_id> <timestamp>
	err = r.client.HSet(ctx, timestampsKey, loc.DriverID.String(), loc.Timestamp.Unix()).Err()
	if err != nil {
		// We don't fail the whole operation if just the timestamp fails,
		// but in a real app we might log it.
		return fmt.Errorf("failed to update driver timestamp in redis: %w", err)
	}

	return nil
}

// GetNearbyDrivers queries the Redis spatial index and fetches exact timestamps.
func (r *locationRepository) GetNearbyDrivers(ctx context.Context, lat, lon float64, radiusKm float64) ([]domain.DriverLocation, error) {
	query := &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Longitude:  lon,
			Latitude:   lat,
			Radius:     radiusKm,
			RadiusUnit: "km",
		},
		WithCoord: true,
	}

	geoLocations, err := r.client.GeoSearchLocation(ctx, locationsKey, query).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to search nearby drivers in redis: %w", err)
	}

	if len(geoLocations) == 0 {
		return []domain.DriverLocation{}, nil
	}

	driverIDs := make([]string, len(geoLocations))
	for i, g := range geoLocations {
		driverIDs[i] = g.Name
	}

	timestamps, err := r.client.HMGet(ctx, timestampsKey, driverIDs...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch timestamps from redis: %w", err)
	}

	drivers := make([]domain.DriverLocation, 0, len(geoLocations))

	for i, geoLoc := range geoLocations {
		driverID, err := uuid.Parse(geoLoc.Name)
		if err != nil {
			continue
		}

		// Redis HMGet returns nil if a key is missing.
		// We handle it gracefully by falling back to current time if timestamp is not found.
		var ts time.Time
		if i < len(timestamps) && timestamps[i] != nil {
			if unixStr, ok := timestamps[i].(string); ok {
				// We store it as Unix string, so we need to parse it back
				var unixTime int64
				fmt.Sscanf(unixStr, "%d", &unixTime)
				ts = time.Unix(unixTime, 0).UTC()
			}
		} else {
			ts = time.Now().UTC()
		}

		drivers = append(drivers, domain.DriverLocation{
			DriverID:  driverID,
			Longitude: geoLoc.Longitude,
			Latitude:  geoLoc.Latitude,
			Timestamp: ts,
		})
	}

	return drivers, nil
}
