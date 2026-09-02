-- name: CreateUser :one
INSERT INTO users (name, email, phone, password_hash, role
) VALUES (
  sqlc.arg(name),
  sqlc.arg(email),
  sqlc.arg(phone),
  sqlc.arg(password_hash),
  sqlc.arg(role)
) RETURNING id, name, email, phone, role, password_hash, created_at, updated_at;

-- name: GetUserById :one
SELECT
  id,
  name,
  email,
  phone,
  role,
  created_at,
  updated_at
from users
WHERE id = sqlc.arg(id);

-- name: GetUserByEmail :one
SELECT
  id,
  name,
  email,
  phone,
  password_hash,
  role,
  created_at,
  updated_at
from users
WHERE LOWER(email) = LOWER(sqlc.arg(email));
