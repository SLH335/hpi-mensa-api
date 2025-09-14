-- name: ListAttributesByMeal :many
SELECT attributes.*
FROM meal_attributes
LEFT JOIN attributes ON meal_attributes.attribute_slug = attributes.slug
WHERE meal_attributes.meal_id = ?;

-- name: ListMealIDsByAttribute :many
SELECT meal_id
FROM meal_attributes
WHERE attribute_slug = ?;

-- name: CreateMealAttribute :exec
INSERT INTO meal_attributes (meal_id, attribute_slug)
VALUES (?, ?);

-- name: DeleteMealAttributesByMeal :exec
DELETE FROM meal_attributes
WHERE meal_id = ?;

-- name: DeleteMealAttributesByAttribute :exec
DELETE FROM meal_attributes
WHERE attribute_slug = ?;
