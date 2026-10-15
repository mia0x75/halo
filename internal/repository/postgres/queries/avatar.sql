-- name: CreateAvatar :one
INSERT INTO mm_avatars (uuid, url, version) VALUES ($1, $2, $3) RETURNING avatar_id;

-- name: GetAvatarByID :one
SELECT avatar_id, uuid, url, version, update_at, create_at FROM mm_avatars WHERE avatar_id = $1 LIMIT 1;

-- name: GetAvatarByUUID :one
SELECT avatar_id, uuid, url, version, update_at, create_at FROM mm_avatars WHERE uuid = $1 LIMIT 1;

-- name: UpdateAvatar :exec
UPDATE mm_avatars SET uuid = $1, url = $2, version = $3, update_at = NOW() WHERE avatar_id = $4;

-- name: DeleteAvatarByID :exec
DELETE FROM mm_avatars WHERE avatar_id = $1;

-- name: DeleteAvatarByUUID :exec
DELETE FROM mm_avatars WHERE uuid = $1;

-- name: ListAvatars :many
SELECT avatar_id, uuid, url, version, update_at, create_at FROM mm_avatars WHERE avatar_id < $1 ORDER BY avatar_id DESC LIMIT $2;

