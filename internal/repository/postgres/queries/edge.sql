-- name: CreateEdge :one
INSERT INTO mm_edges (uuid, "type", ancestor_id, descendant_id, version) VALUES ($1, $2, $3, $4, $5) RETURNING edge_id;

-- name: GetEdgeByID :one
SELECT edge_id, uuid, "type", ancestor_id, descendant_id, version, update_at, create_at FROM mm_edges WHERE edge_id = $1 LIMIT 1;

-- name: GetEdgeByUUID :one
SELECT edge_id, uuid, "type", ancestor_id, descendant_id, version, update_at, create_at FROM mm_edges WHERE uuid = $1 LIMIT 1;

-- name: UpdateEdge :exec
UPDATE mm_edges SET uuid = $1, "type" = $2, ancestor_id = $3, descendant_id = $4, version = $5, update_at = NOW() WHERE edge_id = $6;

-- name: DeleteEdgeByID :exec
DELETE FROM mm_edges WHERE edge_id = $1;

-- name: DeleteEdgeByUUID :exec
DELETE FROM mm_edges WHERE uuid = $1;

-- name: ListEdges :many
SELECT edge_id, uuid, "type", ancestor_id, descendant_id, version, update_at, create_at FROM mm_edges WHERE edge_id < $1 ORDER BY edge_id DESC LIMIT $2;

-- name: DeleteEdgesByTypeAncestor :exec
DELETE FROM mm_edges WHERE "type" = $1 AND ancestor_id = $2;

