-- name: CreateStatistic :exec
INSERT INTO mm_statistics ("group", "key", uuid, value, version) VALUES ($1, $2, $3, $4, $5);

-- name: GetStatisticByUUID :one
SELECT "group", "key", uuid, value, version, update_at, create_at FROM mm_statistics WHERE uuid = $1 LIMIT 1;

-- name: GetStatisticByKey :one
SELECT "group", "key", uuid, value, version, update_at, create_at FROM mm_statistics WHERE "group" = $1 AND "key" = $2 LIMIT 1;

-- name: UpdateStatistic :exec
UPDATE mm_statistics SET value = $1, version = $2, update_at = NOW() WHERE uuid = $3;

-- name: DeleteStatisticByUUID :exec
DELETE FROM mm_statistics WHERE uuid = $1;

-- name: ListStatistics :many
SELECT "group", "key", uuid, value, version, update_at, create_at FROM mm_statistics WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;

-- name: UpsertStatistic :exec
INSERT INTO mm_statistics ("group", "key", uuid, value, version)
    VALUES ($1, $2, $3, $4, 0)
ON CONFLICT ("group", "key") DO UPDATE SET
    value = mm_statistics.value + EXCLUDED.value,
    update_at = NOW();

