-- name: GetMealByID :one
SELECT meals.*, locations.name_de as location_name_de, locations.name_en as location_name_en
FROM meals
JOIN locations ON meals.location_slug = locations.slug
WHERE id = ?;

-- name: ListMealsByMenu :many
SELECT meals.*, locations.name_de as location_name_de, locations.name_en as location_name_en
FROM meals
JOIN locations ON meals.location_slug = locations.slug
WHERE menu_slug = ?;

-- name: CreateMeal :exec
INSERT INTO meals (id, name_de, name_en, category_de, category_en, date, menu_slug, location_slug)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateMealByID :exec
UPDATE meals
SET name_de = ?,
    name_en = ?,
    category_de = ?,
    category_en = ?,
    date = ?,
    menu_slug = ?,
    location_slug = ?
WHERE id = ?;

-- name: DeleteMealByID :exec
DELETE FROM meals
WHERE id = ?;
