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
