-- name: CreateMessage :exec
INSERT INTO mm_messages (uuid, user_id, category, title, content, link) VALUES (?, ?, ?, ?, ?, ?);

-- name: GetMessageByUUID :one
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages WHERE uuid = ? LIMIT 1;

-- name: ListMessagesByUser :many
SELECT message_id, uuid, user_id, category, title, content, link, read_at, create_at FROM mm_messages
WHERE user_id = sqlc.arg(user_id)
  AND message_id < sqlc.arg(before_id)
  AND (read_at IS NULL OR sqlc.arg(unread_only) = 0)
  AND (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(keyword) = '' OR title LIKE sqlc.arg(keyword) OR content LIKE sqlc.arg(keyword))
ORDER BY message_id DESC
LIMIT ?;

-- name: CountMessagesByUser :one
SELECT COUNT(*) FROM mm_messages
WHERE user_id = sqlc.arg(user_id)
  AND (read_at IS NULL OR sqlc.arg(unread_only) = 0)
  AND (sqlc.arg(category) = '' OR category = sqlc.arg(category))
  AND (sqlc.arg(keyword) = '' OR title LIKE sqlc.arg(keyword) OR content LIKE sqlc.arg(keyword));

-- name: CountUnreadMessages :one
SELECT COUNT(*) FROM mm_messages WHERE user_id = ? AND read_at IS NULL;

-- name: MarkMessageRead :exec
UPDATE mm_messages SET read_at = CURRENT_TIMESTAMP WHERE uuid = ? AND user_id = ? AND read_at IS NULL;

-- name: MarkAllMessagesRead :execrows
UPDATE mm_messages SET read_at = CURRENT_TIMESTAMP WHERE user_id = ? AND read_at IS NULL;

-- name: DeleteMessageByUUID :execrows
DELETE FROM mm_messages WHERE uuid = ? AND user_id = ?;

-- name: DeleteMessagesBefore :execrows
DELETE FROM mm_messages WHERE create_at < ?;
