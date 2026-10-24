-- name: CreateUser :one
INSERT INTO mm_users (uuid, email, password, status, name, phone, avatar_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING user_id;

-- name: GetUserByID :one
SELECT user_id, uuid, email, password, status, name, phone, avatar_id, version, update_at, create_at FROM mm_users WHERE user_id = $1 LIMIT 1;

-- name: GetUserByUUID :one
SELECT user_id, uuid, email, password, status, name, phone, avatar_id, version, update_at, create_at FROM mm_users WHERE uuid = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE mm_users SET uuid = $1, email = $2, password = $3, status = $4, name = $5, phone = $6, avatar_id = $7, version = $8, update_at = NOW() WHERE user_id = $9;

-- name: DeleteUserByID :exec
DELETE FROM mm_users WHERE user_id = $1;

-- name: DeleteUserByUUID :exec
DELETE FROM mm_users WHERE uuid = $1;

-- name: ListUsers :many
SELECT user_id, uuid, email, password, status, name, phone, avatar_id, version, update_at, create_at FROM mm_users WHERE user_id < $1 ORDER BY user_id DESC LIMIT $2;


-- name: CountUsers :one
SELECT COUNT(*) FROM mm_users;
