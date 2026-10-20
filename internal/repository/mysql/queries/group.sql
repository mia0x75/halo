-- name: CreateGroup :exec
INSERT INTO mm_groups (uuid, name, description, status, version) VALUES (?, ?, ?, ?, ?);

-- name: GetGroupByUUID :one
SELECT group_id, uuid, name, description, status, version, update_at, create_at FROM mm_groups WHERE uuid = ? LIMIT 1;

-- name: GetGroupByName :one
SELECT group_id, uuid, name, description, status, version, update_at, create_at FROM mm_groups WHERE name = ? LIMIT 1;

-- name: UpdateGroup :exec
UPDATE mm_groups SET name = ?, description = ?, status = ?, version = ?, update_at = CURRENT_TIMESTAMP WHERE uuid = ?;

-- name: DeleteGroupByUUID :exec
DELETE FROM mm_groups WHERE uuid = ?;

-- name: ListGroups :many
SELECT group_id, uuid, name, description, status, version, update_at, create_at FROM mm_groups WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;
