-- name: CreateStatistic :exec
INSERT INTO mm_statistics (`group`, `key`, uuid, value, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetStatisticByUUID :one
SELECT `group`, `key`, uuid, value, version, update_at, create_at FROM mm_statistics WHERE uuid = ? LIMIT 1;

-- name: GetStatisticByKey :one
SELECT `group`, `key`, uuid, value, version, update_at, create_at FROM mm_statistics WHERE `group` = ? AND `key` = ? LIMIT 1;

-- name: UpdateStatistic :exec
UPDATE mm_statistics SET value = ?, version = ?, update_at = ? WHERE uuid = ?;

-- name: DeleteStatisticByUUID :exec
DELETE FROM mm_statistics WHERE uuid = ?;

-- name: ListStatistics :many
SELECT `group`, `key`, uuid, value, version, update_at, create_at FROM mm_statistics WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;
