-- name: CreateMessage :exec
INSERT INTO mm_messages (uuid, user_id, category, title, content, link) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetMessageByUUID :one
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE uuid = ? LIMIT 1;

-- name: ListMessagesByUser :many
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE user_id = ? ORDER BY create_at DESC, message_id DESC LIMIT ?;

-- name: ListUnreadMessagesByUser :many
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE user_id = ? AND read_at IS NULL ORDER BY create_at DESC, message_id DESC LIMIT ?;

-- name: CountUnreadMessages :one
SELECT COUNT(*) FROM mm_messages WHERE user_id = ? AND read_at IS NULL;

-- name: MarkMessageRead :exec
UPDATE mm_messages SET read_at = CURRENT_TIMESTAMP WHERE uuid = ? AND user_id = ? AND read_at IS NULL;

-- name: MarkAllMessagesRead :execrows
UPDATE mm_messages SET read_at = CURRENT_TIMESTAMP WHERE user_id = ? AND read_at IS NULL;
