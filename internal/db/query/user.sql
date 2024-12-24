-- name: CreateUser :one
INSERT INTO "user" (
    id,
    email,
    password
) VALUES (
    $1, $2, $3
) RETURNING *;
