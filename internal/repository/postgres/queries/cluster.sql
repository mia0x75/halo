-- name: CreateCluster :one
INSERT INTO mm_clusters (uuid, host, alias, ip, port, "user", "type", group_id, password, fingerprint, status, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING cluster_id;

-- name: GetClusterByID :one
SELECT cluster_id, uuid, host, alias, ip, port, "user", "type", group_id, password, fingerprint, status, version, update_at, create_at FROM mm_clusters WHERE cluster_id = $1 LIMIT 1;

-- name: GetClusterByUUID :one
SELECT cluster_id, uuid, host, alias, ip, port, "user", "type", group_id, password, fingerprint, status, version, update_at, create_at FROM mm_clusters WHERE uuid = $1 LIMIT 1;

-- name: UpdateCluster :exec
UPDATE mm_clusters SET uuid = $1, host = $2, alias = $3, ip = $4, port = $5, "user" = $6, "type" = $7, group_id = $8, password = $9, fingerprint = $10, status = $11, version = $12, update_at = NOW() WHERE cluster_id = $13;

-- name: DeleteClusterByID :exec
DELETE FROM mm_clusters WHERE cluster_id = $1;

-- name: DeleteClusterByUUID :exec
DELETE FROM mm_clusters WHERE uuid = $1;

-- name: ListClusters :many
SELECT cluster_id, uuid, host, alias, ip, port, "user", "type", group_id, password, fingerprint, status, version, update_at, create_at FROM mm_clusters WHERE cluster_id < $1 ORDER BY cluster_id DESC LIMIT $2;

-- name: CountClustersByGroup :one
SELECT COUNT(*) FROM mm_clusters WHERE group_id = $1;
