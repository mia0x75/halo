-- name: CreateCron :execlastid
INSERT INTO mm_crons (uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, recurrent, hash, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetCronByID :one
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE cron_id = ? LIMIT 1;

-- name: GetCronByUUID :one
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE uuid = ? LIMIT 1;

-- name: UpdateCron :exec
UPDATE mm_crons SET uuid = ?, status = ?, name = ?, cmd = ?, params = ?, last_run = ?, next_run = ?, `interval` = ?, duration = ?, recurrent = ?, hash = ?, version = ?, update_at = ? WHERE cron_id = ?;

-- name: DeleteCronByID :exec
DELETE FROM mm_crons WHERE cron_id = ?;

-- name: DeleteCronByUUID :exec
DELETE FROM mm_crons WHERE uuid = ?;

-- name: ListCrons :many
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, `interval`, duration, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE cron_id < ? ORDER BY cron_id DESC LIMIT ?;
