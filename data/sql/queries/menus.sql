-- name: GetMenuBySlug :one
SELECT *
FROM menus
WHERE slug = ? LIMIT 1;

-- name: ListMenus :many
SELECT *
FROM menus
ORDER BY date;

-- name: ListMenusByDate :many
SELECT *
FROM menus
WHERE DATE(date) = DATE(?);

-- name: ListMenusByLocation :many
SELECT *
FROM menus
WHERE location_slug = ?;

-- name: ListMenusByDateAndLocation :many
SELECT *
FROM menus
WHERE DATE(date) = DATE(?)
    AND location_slug = ?;

-- name: CreateMenu :one
INSERT INTO menus (
    slug, date, provider, location_slug
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: DeleteMenu :exec
DELETE FROM menus
WHERE slug = ?;
