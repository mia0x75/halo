-- name: CreateAvatar :execlastid
INSERT INTO mm_avatars (uuid, url, version) VALUES (?, ?, ?);

-- name: GetAvatarByID :one
SELECT avatar_id, uuid, url, version, update_at, create_at FROM mm_avatars WHERE avatar_id = ? LIMIT 1;

-- name: GetAvatarByUUID :one
SELECT avatar_id, uuid, url, version, update_at, create_at FROM mm_avatars WHERE uuid = ? LIMIT 1;

-- name: UpdateAvatar :exec
UPDATE mm_avatars SET uuid = ?, url = ?, version = ? WHERE avatar_id = ?;

-- name: DeleteAvatarByID :exec
DELETE FROM mm_avatars WHERE avatar_id = ?;

-- name: DeleteAvatarByUUID :exec
DELETE FROM mm_avatars WHERE uuid = ?;

-- name: ListAvatars :many
SELECT avatar_id, uuid, url, version, update_at, create_at FROM mm_avatars WHERE avatar_id < ? ORDER BY avatar_id DESC LIMIT ?;
