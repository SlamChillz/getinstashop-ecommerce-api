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
