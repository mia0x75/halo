-- name: CreateQuery :execlastid
INSERT INTO mm_queries (uuid, `type`, cluster_id, `database`, content, plan, user_id, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetQueryByID :one
SELECT query_id, uuid, `type`, cluster_id, `database`, content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE query_id = ? LIMIT 1;

-- name: GetQueryByUUID :one
SELECT query_id, uuid, `type`, cluster_id, `database`, content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE uuid = ? LIMIT 1;

-- name: UpdateQuery :exec
UPDATE mm_queries SET uuid = ?, `type` = ?, cluster_id = ?, `database` = ?, content = ?, plan = ?, user_id = ?, version = ?, update_at = ? WHERE query_id = ?;

-- name: DeleteQueryByID :exec
DELETE FROM mm_queries WHERE query_id = ?;

-- name: DeleteQueryByUUID :exec
DELETE FROM mm_queries WHERE uuid = ?;

-- name: ListQueries :many
SELECT query_id, uuid, `type`, cluster_id, `database`, content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE query_id < ? ORDER BY query_id DESC LIMIT ?;
