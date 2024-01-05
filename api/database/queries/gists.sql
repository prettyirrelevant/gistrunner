-- name: CreateGist :one
INSERT INTO gists (id, hash, language, result) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListGists :many
SELECT * FROM gists
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetGist :one
SELECT * FROM gists
WHERE hash = $1 LIMIT 1;

-- name: CountGists :one
SELECT COUNT(*) FROM gists;
