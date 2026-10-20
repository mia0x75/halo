-- name: CreateCron :execlastid
INSERT INTO mm_crons (uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, recurrent, hash, version, user_id, cluster_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetCronByID :one
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE cron_id = ? LIMIT 1;

-- name: GetCronByUUID :one
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE uuid = ? LIMIT 1;

-- name: UpdateCron :exec
UPDATE mm_crons SET uuid = ?, status = ?, name = ?, cmd = ?, params = ?, last_run = ?, next_run = ?, `interval` = ?, duration = ?, recurrent = ?, hash = ?, version = ?, user_id = ?, cluster_id = ? WHERE cron_id = ?;

-- name: DeleteCronByID :exec
DELETE FROM mm_crons WHERE cron_id = ?;

-- name: DeleteCronByUUID :exec
DELETE FROM mm_crons WHERE uuid = ?;

-- name: ListCrons :many
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE cron_id < ? ORDER BY cron_id DESC LIMIT ?;

-- name: ListSchedulableCrons :many
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE recurrent = 1 OR (recurrent <> 1 AND status IN ('P', 'H'));
