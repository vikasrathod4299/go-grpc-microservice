-- +goose Up
CREATE TYPE user_role AS ENUM (
    'rider',
    'driver'
);

ALTER TABLE users
ADD COLUMN role user_role NOT NULL DEFAULT 'rider';

-- +goose Down
ALTER TABLE users
DROP COLUMN IF EXISTS role;

DROP TYPE IF EXISTS user_role;
