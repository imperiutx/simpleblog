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

-- name: UpdateUser :one
UPDATE users
SET 
  username = COALESCE(sqlc.narg(username), username),
  email = COALESCE(sqlc.narg(email), email), 
  password = COALESCE(sqlc.narg(password), password), 
  role = COALESCE(sqlc.narg(role), role), 
  status = COALESCE(sqlc.narg(status), status)
WHERE id = sqlc.arg(id)
RETURNING id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
