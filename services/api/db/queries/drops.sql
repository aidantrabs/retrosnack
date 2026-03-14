-- name: ListDrops :many
SELECT id, name, slug, description, instagram_url, released_at, created_at
FROM drops
ORDER BY released_at DESC NULLS LAST, created_at DESC;

-- name: GetDropBySlug :one
SELECT id, name, slug, description, instagram_url, released_at, created_at
FROM drops
WHERE slug = $1;

-- name: GetDropByID :one
SELECT id, name, slug, description, instagram_url, released_at, created_at
FROM drops
WHERE id = $1;

-- name: CreateDrop :one
INSERT INTO drops (name, slug, description, instagram_url)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateDrop :one
UPDATE drops SET
    name          = COALESCE($2, name),
    description   = COALESCE($3, description),
    instagram_url = COALESCE($4, instagram_url)
WHERE id = $1
RETURNING *;

-- name: DeleteDrop :exec
DELETE FROM drops WHERE id = $1;
