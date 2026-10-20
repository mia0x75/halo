-- name: CreateTask :exec
INSERT INTO mm_tasks (uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, error, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTaskByUUID :one
SELECT task_id, uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, started_at, finished_at, error, version, update_at, create_at FROM mm_tasks WHERE uuid = ? LIMIT 1;

-- name: UpdateTask :exec
UPDATE mm_tasks SET name = ?, type = ?, payload = ?, status = ?, priority = ?, attempt = ?, max_attempt = ?, cron_id = ?, scheduled_at = ?, started_at = ?, finished_at = ?, error = ?, version = ?, update_at = CURRENT_TIMESTAMP WHERE uuid = ?;

-- name: DeleteTaskByUUID :exec
DELETE FROM mm_tasks WHERE uuid = ?;

-- name: ListTasks :many
SELECT task_id, uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, started_at, finished_at, error, version, update_at, create_at FROM mm_tasks WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;

-- name: ListPendingTasks :many
SELECT task_id, uuid, name, type, payload, status, priority, attempt, max_attempt, cron_id, scheduled_at, started_at, finished_at, error, version, update_at, create_at FROM mm_tasks WHERE status = 0 ORDER BY priority, scheduled_at LIMIT ?;

-- name: CountTasksByStatus :many
SELECT status, COUNT(*) AS count FROM mm_tasks GROUP BY status ORDER BY status;
