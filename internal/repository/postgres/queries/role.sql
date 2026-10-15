-- name: CreateRole :one
INSERT INTO mm_roles (uuid, name, description, version) VALUES ($1, $2, $3, $4) RETURNING role_id;

-- name: GetRoleByID :one
SELECT role_id, uuid, name, description, version, update_at, create_at FROM mm_roles WHERE role_id = $1 LIMIT 1;

-- name: GetRoleByUUID :one
SELECT role_id, uuid, name, description, version, update_at, create_at FROM mm_roles WHERE uuid = $1 LIMIT 1;

-- name: UpdateRole :exec
UPDATE mm_roles SET uuid = $1, name = $2, description = $3, version = $4, update_at = NOW() WHERE role_id = $5;

-- name: DeleteRoleByID :exec
DELETE FROM mm_roles WHERE role_id = $1;

-- name: DeleteRoleByUUID :exec
DELETE FROM mm_roles WHERE uuid = $1;

-- name: ListRoles :many
SELECT role_id, uuid, name, description, version, update_at, create_at FROM mm_roles WHERE role_id < $1 ORDER BY role_id DESC LIMIT $2;

