-- name: CreateUser :one
INSERT INTO users (name, email, phone, password
) VALUES (
$1, $2, $3, $4
) RETURNING id, name, email, phone, created_at, updated_at;

-- name GetUserById :one
SELECT
  id,
  name,
  email,
  phone,
  created_at,
  updated_at
from users
WHERE id = $1;
