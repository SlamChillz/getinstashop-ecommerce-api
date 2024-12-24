-- name: CreateUser :one
INSERT INTO "user" (
    id,
    email,
    password
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetUserById :one
SELECT id, email, password, admin
FROM "user"
WHERE email = $1 LIMIT 1;
