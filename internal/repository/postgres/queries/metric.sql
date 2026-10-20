-- name: CreateMetric :exec
INSERT INTO mm_metrics (uuid, name, scope, scope_id, type, bucket, value, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetMetricByUUID :one
SELECT metric_id, uuid, name, scope, scope_id, type, bucket, value, version, update_at, create_at FROM mm_metrics WHERE uuid = $1 LIMIT 1;

-- name: GetMetricByKey :one
SELECT metric_id, uuid, name, scope, scope_id, type, bucket, value, version, update_at, create_at FROM mm_metrics WHERE name = $1 AND scope = $2 AND scope_id = $3 AND bucket = $4 LIMIT 1;

-- name: UpdateMetric :exec
UPDATE mm_metrics SET value = $1, version = $2, update_at = NOW() WHERE uuid = $3;

-- name: DeleteMetricByUUID :exec
DELETE FROM mm_metrics WHERE uuid = $1;

-- name: ListMetrics :many
SELECT metric_id, uuid, name, scope, scope_id, type, bucket, value, version, update_at, create_at FROM mm_metrics WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;

-- name: IncrementMetric :exec
-- Counter semantics: the new value is the old value plus the delta, so
-- concurrent writers never lose an update.
INSERT INTO mm_metrics (uuid, name, scope, scope_id, type, bucket, value, version)
    VALUES ($1, $2, $3, $4, 0, $5, $6, 0)
ON CONFLICT (name, scope, scope_id, bucket) DO UPDATE SET
    value = mm_metrics.value + EXCLUDED.value,
    update_at = NOW();

-- name: SetGaugeMetric :exec
-- Gauge semantics: the value is the new value, not an increment.
INSERT INTO mm_metrics (uuid, name, scope, scope_id, type, bucket, value, version)
    VALUES ($1, $2, $3, $4, 1, $5, $6, 0)
ON CONFLICT (name, scope, scope_id, bucket) DO UPDATE SET
    value = EXCLUDED.value,
    update_at = NOW();
