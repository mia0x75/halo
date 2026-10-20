-- name: CreateCluster :execlastid
INSERT INTO mm_clusters (uuid, host, alias, ip, port, `user`, `type`, group_id, password, fingerprint, status, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetClusterByID :one
SELECT cluster_id, uuid, host, alias, ip, port, `user`, `type`, group_id, password, fingerprint, status, version, update_at, create_at FROM mm_clusters WHERE cluster_id = ? LIMIT 1;

-- name: GetClusterByUUID :one
SELECT cluster_id, uuid, host, alias, ip, port, `user`, `type`, group_id, password, fingerprint, status, version, update_at, create_at FROM mm_clusters WHERE uuid = ? LIMIT 1;

-- name: UpdateCluster :exec
UPDATE mm_clusters SET uuid = ?, host = ?, alias = ?, ip = ?, port = ?, `user` = ?, `type` = ?, group_id = ?, password = ?, fingerprint = ?, status = ?, version = ? WHERE cluster_id = ?;

-- name: DeleteClusterByID :exec
DELETE FROM mm_clusters WHERE cluster_id = ?;

-- name: DeleteClusterByUUID :exec
DELETE FROM mm_clusters WHERE uuid = ?;

-- name: ListClusters :many
SELECT cluster_id, uuid, host, alias, ip, port, `user`, `type`, group_id, password, fingerprint, status, version, update_at, create_at FROM mm_clusters WHERE cluster_id < ? ORDER BY cluster_id DESC LIMIT ?;
