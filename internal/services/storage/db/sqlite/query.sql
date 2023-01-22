-- name: CreateAsset :one
INSERT INTO assets (
       name,
       path
) VALUES (
  ?, ?
)
RETURNING *;

-- name: GetAsset :one
SELECT *
FROM assets
WHERE name = ?;

-- name: GetAssets :many
SELECT *
FROM assets;
