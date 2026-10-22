-- name: CreateOption :exec
INSERT INTO mm_options (name, uuid, parameter, description, element, writable, version) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetOptionByUUID :one
SELECT name, uuid, parameter, description, element, writable, group_name, sort_order, depends_on, min_value, max_value, version, update_at, create_at FROM mm_options WHERE uuid = $1 LIMIT 1;

-- name: GetOptionByName :one
SELECT name, uuid, parameter, description, element, writable, group_name, sort_order, depends_on, min_value, max_value, version, update_at, create_at FROM mm_options WHERE name = $1 LIMIT 1;

-- name: UpdateOption :exec
UPDATE mm_options SET parameter = $1, description = $2, element = $3, writable = $4, version = $5, update_at = NOW() WHERE uuid = $6;

-- name: DeleteOptionByUUID :exec
DELETE FROM mm_options WHERE uuid = $1;

-- name: ListOptions :many
SELECT name, uuid, parameter, description, element, writable, group_name, sort_order, depends_on, min_value, max_value, version, update_at, create_at FROM mm_options ORDER BY sort_order, name LIMIT $1;
