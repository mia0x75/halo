-- name: CreateMessage :exec
INSERT INTO mm_messages (uuid, user_id, category, title, content, link) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetMessageByUUID :one
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE uuid = $1 LIMIT 1;

-- name: ListMessagesByUser :many
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE user_id = $1 ORDER BY create_at DESC, message_id DESC LIMIT $2;

-- name: ListUnreadMessagesByUser :many
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE user_id = $1 AND read_at IS NULL ORDER BY create_at DESC, message_id DESC LIMIT $2;

-- name: CountUnreadMessages :one
SELECT COUNT(*) FROM mm_messages WHERE user_id = $1 AND read_at IS NULL;

-- name: MarkMessageRead :exec
UPDATE mm_messages SET read_at = NOW() WHERE uuid = $1 AND user_id = $2 AND read_at IS NULL;

-- name: MarkAllMessagesRead :execrows
UPDATE mm_messages SET read_at = NOW() WHERE user_id = $1 AND read_at IS NULL;
