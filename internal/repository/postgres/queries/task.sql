-- name: CreateTask :exec
INSERT INTO mm_tasks (uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, error, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);

-- name: GetTaskByUUID :one
SELECT task_id, uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, started_at, finished_at, error, version, update_at, create_at FROM mm_tasks WHERE uuid = $1 LIMIT 1;

-- name: UpdateTask :exec
UPDATE mm_tasks SET name = $1, type = $2, payload = $3, status = $4, priority = $5, attempt = $6, max_attempt = $7, cron_id = $8, scheduled_at = $9, started_at = $10, finished_at = $11, error = $12, version = $13, update_at = NOW() WHERE uuid = $14;

-- name: DeleteTaskByUUID :exec
DELETE FROM mm_tasks WHERE uuid = $1;

-- name: ListTasks :many
SELECT task_id, uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, started_at, finished_at, error, version, update_at, create_at FROM mm_tasks WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;

-- name: ListPendingTasks :many
SELECT task_id, uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, started_at, finished_at, error, version, update_at, create_at FROM mm_tasks WHERE status = 0 ORDER BY priority, scheduled_at LIMIT $1;

-- name: CountTasksByStatus :many
SELECT status, COUNT(*) AS count FROM mm_tasks GROUP BY status ORDER BY status;
