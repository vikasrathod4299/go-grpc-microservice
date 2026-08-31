package repository

/*
================================================================================
FILE: internal/location/repository/redis_geo.go
================================================================================

PURPOSE:
Database access layer interacting directly with Redis using `go-redis/v9`.
Executes high-speed in-memory spatial commands: `GEOADD`, `GEOSEARCH`, `ZREM`.

LEARNING GO CONCEPTS:
- Using `github.com/redis/go-redis/v9`.
- Managing Redis keys (e.g., `"drivers:available"`).
- Processing Redis GeoLocation query results.

WHAT YOU NEED TO IMPLEMENT HERE:
1. `RedisGeoRepository` struct wrapping `*redis.Client`.

2. `SaveLocation(ctx context.Context, driverID string, lat, lng float64) error`
   - Execute: `r.client.GeoAdd(ctx, "drivers:available", &redis.GeoLocation{ Name: driverID, Latitude: lat, Longitude: lng })`

3. `FindNearby(ctx context.Context, lat, lng, radiusKM float64, limit int) ([]service.NearbyDriver, error)`
   - Execute `GEOSEARCH drivers:available FROMLONLAT lng lat BYRADIUS radiusKM km WITHDIST WITHCOORD COUNT limit ASC`
   - Map Redis locations to `service.NearbyDriver` slices.

4. `Remove(ctx context.Context, driverID string) error`
   - Execute: `r.client.ZRem(ctx, "drivers:available", driverID)`
================================================================================
*/

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/vikasrathod4299/microservice/internal/location/service"
)

const driversKey = "drivers:available"

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

type RedisGeoRepository struct {
	client *redis.Client
}

func NewRedisGeoRepository(client *redis.Client) *RedisGeoRepository {
	return &RedisGeoRepository{client: client}
}

func (r *RedisGeoRepository) SaveLocation(ctx context.Context, driverID string, lat, lng float64) error {
	return r.client.GeoAdd(ctx, driversKey, &redis.GeoLocation{
		Name:      driverID,
		Latitude:  lat,
		Longitude: lng,
	}).Err()
}

func (r *RedisGeoRepository) FindNearby(ctx context.Context, lat, lng, rediusKM float64, limit int) ([]service.NearbyDriver, error) {
	results, err := r.client.GeoSearchLocation(ctx, driversKey, &redis.GeoSearchLocationQuery{
		GeoSearchQuery: redis.GeoSearchQuery{
			Latitude:   lat,
			Longitude:  lng,
			Radius:     rediusKM,
			RadiusUnit: "km",
			Sort:       "ASC",
			Count:      limit,
			CountAny:   false,
		},
		WithDist:  true,
		WithCoord: true,
	}).Result()
	if err != nil {
		return nil, err
	}
	drivers := make([]service.NearbyDriver, 0, len(results))

	for _, loc := range results {
		drivers = append(drivers, service.NearbyDriver{
			Latitude:   loc.Latitude,
			Longitude:  loc.Longitude,
			DistanceKM: loc.Dist,
			DriverID:   loc.Name,
		})
	}

	return drivers, nil
}

func (r *RedisGeoRepository) Remove(ctx context.Context, diriverID string) error {
	return r.client.ZRem(ctx, driversKey, diriverID).Err()
}
