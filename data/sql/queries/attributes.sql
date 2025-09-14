-- name: GetAttributeBySlug :one
SELECT *
FROM attributes
WHERE slug = ?;

-- name: ListAttributesForMeals :many
SELECT ma.meal_id, a.*
FROM meal_attributes ma
JOIN attributes a ON ma.attribute_slug = a.slug
WHERE ma.meal_id IN (sqlc.slice('meal_ids'));

-- name: CreateAttribute :exec
INSERT INTO attributes (slug, type, name_de, name_en, short_de, short_en)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateAttribute :exec
UPDATE attributes
SET type = ?, name_de = ?, name_en = ?, short_de = ?, short_en = ?
WHERE slug = ?;

-- name: DeleteAttributeBySlug :exec
DELETE FROM attributes
WHERE slug = ?;
