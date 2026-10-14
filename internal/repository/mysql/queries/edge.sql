-- name: CreateEdge :execlastid
INSERT INTO mm_edges (uuid, `type`, ancestor_id, descendant_id, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetEdgeByID :one
SELECT edge_id, uuid, `type`, ancestor_id, descendant_id, version, update_at, create_at FROM mm_edges WHERE edge_id = ? LIMIT 1;

-- name: GetEdgeByUUID :one
SELECT edge_id, uuid, `type`, ancestor_id, descendant_id, version, update_at, create_at FROM mm_edges WHERE uuid = ? LIMIT 1;

-- name: UpdateEdge :exec
UPDATE mm_edges SET uuid = ?, `type` = ?, ancestor_id = ?, descendant_id = ?, version = ?, update_at = ? WHERE edge_id = ?;

-- name: DeleteEdgeByID :exec
DELETE FROM mm_edges WHERE edge_id = ?;

-- name: DeleteEdgeByUUID :exec
DELETE FROM mm_edges WHERE uuid = ?;

-- name: ListEdges :many
SELECT edge_id, uuid, `type`, ancestor_id, descendant_id, version, update_at, create_at FROM mm_edges WHERE edge_id < ? ORDER BY edge_id DESC LIMIT ?;
