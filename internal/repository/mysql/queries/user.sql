-- name: CreateUser :execlastid
INSERT INTO mm_users (uuid, email, password, status, name, phone, avatar_id, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetUserByID :one
SELECT user_id, uuid, email, password, status, name, phone, avatar_id, version, update_at, create_at FROM mm_users WHERE user_id = ? LIMIT 1;

-- name: GetUserByUUID :one
SELECT user_id, uuid, email, password, status, name, phone, avatar_id, version, update_at, create_at FROM mm_users WHERE uuid = ? LIMIT 1;

-- name: UpdateUser :exec
UPDATE mm_users SET uuid = ?, email = ?, password = ?, status = ?, name = ?, phone = ?, avatar_id = ?, version = ? WHERE user_id = ?;

-- name: DeleteUserByID :exec
DELETE FROM mm_users WHERE user_id = ?;

-- name: DeleteUserByUUID :exec
DELETE FROM mm_users WHERE uuid = ?;

-- name: ListUsers :many
SELECT user_id, uuid, email, password, status, name, phone, avatar_id, version, update_at, create_at FROM mm_users WHERE user_id < ? ORDER BY user_id DESC LIMIT ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM mm_users;
