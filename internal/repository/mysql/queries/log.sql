-- name: CreateLog :execlastid
INSERT INTO mm_logs (uuid, user_id, operation, version) VALUES (?, ?, ?, ?);

-- name: GetLogByID :one
SELECT log_id, uuid, user_id, operation, version, create_at FROM mm_logs WHERE log_id = ? LIMIT 1;

-- name: GetLogByUUID :one
SELECT log_id, uuid, user_id, operation, version, create_at FROM mm_logs WHERE uuid = ? LIMIT 1;

-- name: UpdateLog :exec
UPDATE mm_logs SET uuid = ?, user_id = ?, operation = ?, version = ? WHERE log_id = ?;

-- name: DeleteLogByID :exec
DELETE FROM mm_logs WHERE log_id = ?;

-- name: DeleteLogByUUID :exec
DELETE FROM mm_logs WHERE uuid = ?;

-- name: ListLogs :many
SELECT log_id, uuid, user_id, operation, version, create_at FROM mm_logs WHERE log_id < ? ORDER BY log_id DESC LIMIT ?;

-- name: DeleteLogsBefore :execrows
DELETE FROM mm_logs WHERE create_at < ?;
