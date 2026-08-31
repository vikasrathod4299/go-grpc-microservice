-- name: CreateDriver :one
INSERT INTO drivers (
    id, name, phone, vehicle_make, vehicle_model, license_plate, availability, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW(), NOW()
)
RETURNING *;

-- name: GetDriverByID :one
SELECT * FROM drivers
WHERE id = $1 LIMIT 1;

-- name: UpdateDriverAvailability :exec
UPDATE drivers
SET availability = $2, updated_at = NOW()
WHERE id = $1;
