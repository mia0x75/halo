-- name: CreateComment :execlastid
INSERT INTO mm_comments (uuid, content, ticket_id, user_id, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetCommentByID :one
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE comment_id = ? LIMIT 1;

-- name: GetCommentByUUID :one
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE uuid = ? LIMIT 1;

-- name: UpdateComment :exec
UPDATE mm_comments SET uuid = ?, content = ?, ticket_id = ?, user_id = ?, version = ?, update_at = ? WHERE comment_id = ?;

-- name: DeleteCommentByID :exec
DELETE FROM mm_comments WHERE comment_id = ?;

-- name: DeleteCommentByUUID :exec
DELETE FROM mm_comments WHERE uuid = ?;

-- name: ListComments :many
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE comment_id < ? ORDER BY comment_id DESC LIMIT ?;
