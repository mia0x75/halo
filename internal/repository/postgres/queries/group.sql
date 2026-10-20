-- name: CreateGroup :exec
INSERT INTO mm_groups (uuid, name, description, status, version) VALUES ($1, $2, $3, $4, $5);

-- name: GetGroupByUUID :one
SELECT group_id, uuid, name, description, status, version, update_at, create_at FROM mm_groups WHERE uuid = $1 LIMIT 1;

-- name: GetGroupByName :one
SELECT group_id, uuid, name, description, status, version, update_at, create_at FROM mm_groups WHERE name = $1 LIMIT 1;

-- name: UpdateGroup :exec
UPDATE mm_groups SET name = $1, description = $2, status = $3, version = $4, update_at = NOW() WHERE uuid = $5;

-- name: DeleteGroupByUUID :exec
DELETE FROM mm_groups WHERE uuid = $1;

-- name: ListGroups :many
SELECT group_id, uuid, name, description, status, version, update_at, create_at FROM mm_groups WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;
