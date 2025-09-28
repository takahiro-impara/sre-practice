-- name: CreateUser :one
INSERT INTO users (
    email,
    password,
    name
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: CreateUserWithID :one
INSERT INTO users (
    id,
    email,
    password,
    name
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: CheckUserExistsByID :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);

-- name: CheckUserExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByName :one
SELECT * FROM users WHERE name = $1;

-- name: UpdateUser :one
UPDATE users SET email = $1, name = $2, updated_at = NOW() WHERE id = $3 RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2;
