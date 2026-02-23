-- name: CreateDriver :exec
INSERT INTO drivers (id, name, status, created_at, updated_at)
VALUES (@id, @name, @status, @created_at, @updated_at);

-- name: GetDriverByID :one
SELECT id, name, status, created_at, updated_at
FROM drivers
WHERE id = @id LIMIT 1;

-- name: UpdateDriverStatus :exec
UPDATE drivers
SET status = @status, updated_at = @updated_at
WHERE id = @id;

-- name: UpdateDriverLocation :exec
-- We use ST_SetSRID and ST_MakePoint to let Postgres handle the complex spatial data creation.
UPDATE drivers
SET current_location = ST_SetSRID(ST_MakePoint(@lon::float, @lat::float), 4326)::geography, 
    updated_at = @updated_at
WHERE id = @id;