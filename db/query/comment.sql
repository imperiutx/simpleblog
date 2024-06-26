-- name: CreateComment :one
INSERT INTO comments (
  username,
  post_id,
  content
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetCommentById :one
SELECT * FROM comments
WHERE id = $1 LIMIT 1;

-- name: UpdateComment :one
UPDATE comments
SET
  content = COALESCE(sqlc.narg(content), content)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;

-- name: ListCommentsByPostID :many
SELECT *
FROM comments
WHERE post_id = $1
ORDER BY id;

-- name: ListAllComments :many
SELECT *
FROM comments
ORDER BY id DESC
LIMIT $1
OFFSET $2;