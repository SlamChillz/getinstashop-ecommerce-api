-- name: CreateProduct :one
INSERT INTO "product" (
    id,
    name,
    description,
    price,
    stock,
    "createdBy"
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAllProduct :many
SELECT
    id,
    name,
    description,
    price,
    stock,
    "createdBy",
    "createdAt",
    "updatedAt"
FROM "product";

-- name: GetOneProduct :one
SELECT
    id,
    name,
    description,
    price,
    stock,
    "createdBy",
    "createdAt",
    "updatedAt"
FROM "product"
WHERE id = $1
LIMIT 1;

-- name: DeleteOneProduct :exec
DELETE FROM product
WHERE id = $1;

-- name: UpdateOneProduct :one
UPDATE product
SET
    name = sqlc.arg('name'),
    description = sqlc.arg('description'),
    price = sqlc.arg('price'),
    stock = sqlc.arg('stock')
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateProductStock :one
UPDATE product
SET
    stock = stock - $2
WHERE id = $1
RETURNING *;

-- name: GetMultipleProductById :many
SELECT
    id,
    price,
    stock
FROM product
WHERE id = ANY($1::UUID[]);
