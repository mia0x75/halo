-- name: CreateOption :exec
INSERT INTO mm_options (name, uuid, parameter, description, element, writable, version) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetOptionByUUID :one
SELECT name, uuid, parameter, description, element, writable, group_name, sort_order, depends_on, min_value, max_value, version, update_at, create_at FROM mm_options WHERE uuid = ? LIMIT 1;

-- name: GetOptionByName :one
SELECT name, uuid, parameter, description, element, writable, group_name, sort_order, depends_on, min_value, max_value, version, update_at, create_at FROM mm_options WHERE name = ? LIMIT 1;

-- name: UpdateOption :exec
UPDATE mm_options SET parameter = ?, description = ?, element = ?, writable = ?, version = ? WHERE uuid = ?;

-- name: DeleteOptionByUUID :exec
DELETE FROM mm_options WHERE uuid = ?;

-- name: ListOptions :many
SELECT name, uuid, parameter, description, element, writable, group_name, sort_order, depends_on, min_value, max_value, version, update_at, create_at FROM mm_options ORDER BY sort_order, name LIMIT ?;
