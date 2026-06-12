-- name: CreateUser :exec
INSERT INTO users (name, dob)
VALUES (?, ?);

-- name: GetUserByID :one
SELECT id, name, dob
FROM users
WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users
SET name = ?, dob = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?;

-- name: ListUsers :many
SELECT id, name, dob
FROM users
ORDER BY id ASC
LIMIT ?
OFFSET ?;