-- name: GetLocationBySlug :one
SELECT *
FROM locations
WHERE slug = ?;

-- name: ListLocations :many
SELECT *
FROM locations;

-- name: ListLocationsByProvider :many
SELECT *
FROM locations
WHERE provider = ?;

-- name: CreateLocation :exec
INSERT INTO locations (slug, name_de, name_en, provider)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: DeleteLocationBySlug :exec
DELETE FROM locations
WHERE slug = ?;
