-- name: CreateRole :execlastid
INSERT INTO mm_roles (uuid, name, description, version) VALUES (?, ?, ?, ?);

-- name: GetRoleByID :one
SELECT role_id, uuid, name, description, version, update_at, create_at FROM mm_roles WHERE role_id = ? LIMIT 1;

-- name: GetRoleByUUID :one
SELECT role_id, uuid, name, description, version, update_at, create_at FROM mm_roles WHERE uuid = ? LIMIT 1;

-- name: UpdateRole :exec
UPDATE mm_roles SET uuid = ?, name = ?, description = ?, version = ? WHERE role_id = ?;

-- name: DeleteRoleByID :exec
DELETE FROM mm_roles WHERE role_id = ?;

-- name: DeleteRoleByUUID :exec
DELETE FROM mm_roles WHERE uuid = ?;

-- name: ListRoles :many
SELECT role_id, uuid, name, description, version, update_at, create_at FROM mm_roles WHERE role_id < ? ORDER BY role_id DESC LIMIT ?;
