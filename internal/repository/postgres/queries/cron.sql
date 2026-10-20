-- name: CreateCron :one
INSERT INTO mm_crons (uuid, status, name, cmd, params, last_run, next_run, "interval", duration, recurrent, hash, version, user_id, cluster_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING cron_id;

-- name: GetCronByID :one
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, "interval", duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE cron_id = $1 LIMIT 1;

-- name: GetCronByUUID :one
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, "interval", duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE uuid = $1 LIMIT 1;

-- name: UpdateCron :exec
UPDATE mm_crons SET uuid = $1, status = $2, name = $3, cmd = $4, params = $5, last_run = $6, next_run = $7, "interval" = $8, duration = $9, recurrent = $10, hash = $11, version = $12, user_id = $13, cluster_id = $14, update_at = NOW() WHERE cron_id = $15;

-- name: DeleteCronByID :exec
DELETE FROM mm_crons WHERE cron_id = $1;

-- name: DeleteCronByUUID :exec
DELETE FROM mm_crons WHERE uuid = $1;

-- name: ListCrons :many
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, "interval", duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE cron_id < $1 ORDER BY cron_id DESC LIMIT $2;

-- name: ListSchedulableCrons :many
SELECT cron_id, uuid, status, name, cmd, params, last_run, next_run, "interval", duration, user_id, cluster_id, recurrent, hash, version, update_at, create_at FROM mm_crons WHERE recurrent = 1 OR (recurrent <> 1 AND status IN ('P', 'H'));

