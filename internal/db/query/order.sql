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

-- name: GetAllOrderByUserId :many
SELECT * FROM "order"
WHERE "userId" = $1;

-- name: GetAllOrderItem :many
SELECT * FROM "orderItem"
WHERE "orderId" = $1;

-- name: CancelOrder :one
UPDATE "order"
SET
    status = 'CANCELLED'
WHERE id = $1 AND "userId" = $2 AND status = 'PENDING'
RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE "order"
SET
    status = sqlc.arg('status')
WHERE id = sqlc.arg('id')
RETURNING *;
