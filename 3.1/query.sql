-- name: GetUser :one
SELECT * FROM users
WHERE uuid = ? LIMIT 1;

-- name: CreateUser :exec
INSERT INTO users (
  uuid, email, password, first_name, last_name, nickname, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: UpdateUser :exec
UPDATE users SET email = ?, password = ?, first_name = ?, last_name = ?, nickname = ?, updated_at = ? WHERE uuid = ?;

-- name: FindUserByEmailPass :one
SELECT * FROM users
WHERE email = ? and password = ? LIMIT 1;

-- name: FindAllUsers :many
SELECT * FROM users ORDER BY created_at desc LIMIT ? OFFSET ?;

