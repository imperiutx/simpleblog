-- name: CreateComment :one
INSERT INTO comments (
  user_id,
  post_id,
  content
) VALUES (
  $1, $2, $3
) RETURNING id, created_at;

-- name: UpdateComment :one
UPDATE comments
SET
  content = COALESCE(sqlc.narg(content), content)
WHERE id = sqlc.arg(id)
RETURNING id, edited_at;

-- name: DeleteComment :exec
DELETE FROM comments
WHERE id = $1;
