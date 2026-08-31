-- name: CreateTrip :one
INSERT INTO trips (
    id, rider_id, driver_id, status, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
)
RETURNING id, rider_id, driver_id, status, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, created_at, updated_at;

-- name: GetTripByID :one
SELECT
    id,
    rider_id,
    driver_id,
    status,
    pickup_lat,
    pickup_lng,
    dropoff_lat,
    dropoff_lng,
    created_at,
    updated_at
FROM trips
WHERE id = $1;

-- name: UpdateTripStatus :one
UPDATE trips
SET
    status = sqlc.arg(status),
    driver_id = COALESCE(
        NULLIF(sqlc.arg(driver_id)::text, ''),
        driver_id
    ),
    updated_at = NOW()
WHERE id = sqlc.arg(id)
RETURNING
    id,
    rider_id,
    driver_id,
    status,
    pickup_lat,
    pickup_lng,
    dropoff_lat,
    dropoff_lng,
    created_at,
    updated_at;
