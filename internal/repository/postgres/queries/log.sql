-- name: CreateLog :one
INSERT INTO mm_logs (uuid, user_id, operation, version) VALUES ($1, $2, $3, $4) RETURNING log_id;

-- name: GetLogByID :one
SELECT log_id, uuid, user_id, operation, version, create_at FROM mm_logs WHERE log_id = $1 LIMIT 1;

-- name: GetLogByUUID :one
SELECT log_id, uuid, user_id, operation, version, create_at FROM mm_logs WHERE uuid = $1 LIMIT 1;

-- name: UpdateLog :exec
UPDATE mm_logs SET uuid = $1, user_id = $2, operation = $3, version = $4 WHERE log_id = $5;

-- name: DeleteLogByID :exec
DELETE FROM mm_logs WHERE log_id = $1;

-- name: DeleteLogByUUID :exec
DELETE FROM mm_logs WHERE uuid = $1;

-- name: ListLogs :many
SELECT log_id, uuid, user_id, operation, version, create_at FROM mm_logs WHERE log_id < $1 ORDER BY log_id DESC LIMIT $2;

-- name: DeleteLogsBefore :execrows
DELETE FROM mm_logs WHERE create_at < $1;
