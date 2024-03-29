-- name: CreatePost :one
INSERT INTO posts (
  username,
  title,
  content, 
  tags
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetPostById :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1;

-- name: GetPostForUpdate :one
SELECT * FROM posts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListPosts :many
SELECT * FROM posts
ORDER BY id DESC;

-- name: UpdatePost :one
UPDATE posts
SET 
  title = COALESCE(sqlc.narg(title), title),
  content = COALESCE(sqlc.narg(content), content),
  edited_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeletePost :exec
DELETE FROM posts
WHERE id = $1;
