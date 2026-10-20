-- name: CreateMetric :exec
INSERT INTO mm_metrics (uuid, name, scope, scope_id, type, bucket, value, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetMetricByUUID :one
SELECT metric_id, uuid, name, scope, scope_id, type, bucket, value, version, update_at, create_at FROM mm_metrics WHERE uuid = ? LIMIT 1;

-- name: GetMetricByKey :one
SELECT metric_id, uuid, name, scope, scope_id, type, bucket, value, version, update_at, create_at FROM mm_metrics WHERE name = ? AND scope = ? AND scope_id = ? AND bucket = ? LIMIT 1;

-- name: UpdateMetric :exec
UPDATE mm_metrics SET value = ?, version = ? WHERE uuid = ?;

-- name: DeleteMetricByUUID :exec
DELETE FROM mm_metrics WHERE uuid = ?;

-- name: ListMetrics :many
SELECT metric_id, uuid, name, scope, scope_id, type, bucket, value, version, update_at, create_at FROM mm_metrics WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;

-- name: IncrementMetric :exec
-- Counter semantics: the new value is the old value plus the delta, so
-- concurrent writers never lose an update.
INSERT INTO mm_metrics (uuid, name, scope, scope_id, type, bucket, value, version)
    VALUES (?, ?, ?, ?, 0, ?, ?, 0)
ON DUPLICATE KEY UPDATE
    value = value + VALUES(value),
    update_at = CURRENT_TIMESTAMP;

-- name: SetGaugeMetric :exec
-- Gauge semantics: the delta is the new value, not an increment.
INSERT INTO mm_metrics (uuid, name, scope, scope_id, type, bucket, value, version)
    VALUES (?, ?, ?, ?, 1, ?, ?, 0)
ON DUPLICATE KEY UPDATE
    value = VALUES(value),
    update_at = CURRENT_TIMESTAMP;
