-- name: CreateUser :one
INSERT INTO users (
  username,
  password,
  email
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserForUpdate :one
SELECT * FROM users
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: UpdateUser :one
UPDATE users
SET 
  password = COALESCE(sqlc.narg(password), password),
  status = COALESCE(sqlc.narg(status), status)
WHERE id = sqlc.arg(id)
RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
