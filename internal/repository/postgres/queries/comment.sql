-- name: CreateComment :one
INSERT INTO mm_comments (uuid, content, ticket_id, user_id, version) VALUES ($1, $2, $3, $4, $5) RETURNING comment_id;

-- name: GetCommentByID :one
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE comment_id = $1 LIMIT 1;

-- name: GetCommentByUUID :one
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE uuid = $1 LIMIT 1;

-- name: UpdateComment :exec
UPDATE mm_comments SET uuid = $1, content = $2, ticket_id = $3, user_id = $4, version = $5, update_at = NOW() WHERE comment_id = $6;

-- name: DeleteCommentByID :exec
DELETE FROM mm_comments WHERE comment_id = $1;

-- name: DeleteCommentByUUID :exec
DELETE FROM mm_comments WHERE uuid = $1;

-- name: ListComments :many
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE comment_id < $1 ORDER BY comment_id DESC LIMIT $2;

-- name: ListCommentsByTicket :many
SELECT comment_id, uuid, content, ticket_id, user_id, version, update_at, create_at FROM mm_comments WHERE ticket_id = $1 ORDER BY comment_id;

