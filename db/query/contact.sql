-- name: CreateContact :one
INSERT INTO contacts (
  first_name,
  last_name,
  email,
  phone
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetContactById :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetContactById :one
SELECT count(*) FROM contacts
WHERE email = $1;

-- name: GetContactForUpdate :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY id;

-- name: UpdateContact :one
UPDATE contacts
SET 
  first_name = COALESCE(sqlc.narg(first_name), first_name),
  last_name = COALESCE(sqlc.narg(last_name), last_name),
  phone = COALESCE(sqlc.narg(phone), phone),
  updated_at = now()
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = $1;
