-- name: GetPriceById :one
SELECT *
FROM prices
WHERE id = ?;

-- name: ListPricesByMeal :many
SELECT *
FROM prices
WHERE meal_id = ?;

-- name: CreatePrice :exec
INSERT INTO prices (meal_id, type_de, type_en, amount)
VALUES (?, ?, ?, ?);

-- name: UpdatePrice :exec
UPDATE prices
SET amount = ?
WHERE id = ?;

-- name: DeletePriceById :exec
DELETE FROM prices
WHERE id = ?;

-- name: DeletePricesByMeal :exec
DELETE FROM prices
WHERE meal_id = ?;
