-- name: CreateOrder :one
INSERT INTO "order" (
    id,
    "userId",
    total
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetOrderById :one
SELECT * FROM "order"
WHERE id = $1;

-- name: GetAllOrderByUserId :one
SELECT * FROM "order"
WHERE "userId" = $1;

-- name: GetAllOrderItem :many
SELECT * FROM "orderItem"
WHERE "orderId" = $1;
