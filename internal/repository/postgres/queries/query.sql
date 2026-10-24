-- name: CreateQuery :one
INSERT INTO mm_queries (uuid, "type", cluster_id, "database", content, plan, user_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING query_id;

-- name: GetQueryByID :one
SELECT query_id, uuid, "type", cluster_id, "database", content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE query_id = $1 LIMIT 1;

-- name: GetQueryByUUID :one
SELECT query_id, uuid, "type", cluster_id, "database", content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE uuid = $1 LIMIT 1;

-- name: UpdateQuery :exec
UPDATE mm_queries SET uuid = $1, "type" = $2, cluster_id = $3, "database" = $4, content = $5, plan = $6, user_id = $7, version = $8, update_at = NOW() WHERE query_id = $9;

-- name: DeleteQueryByID :exec
DELETE FROM mm_queries WHERE query_id = $1;

-- name: DeleteQueryByUUID :exec
DELETE FROM mm_queries WHERE uuid = $1;

-- name: ListQueries :many
SELECT query_id, uuid, "type", cluster_id, "database", content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE query_id < $1 ORDER BY query_id DESC LIMIT $2;

-- name: ListQueriesByUser :many
SELECT query_id, uuid, "type", cluster_id, "database", content, plan, user_id, version, update_at, create_at FROM mm_queries WHERE user_id = $1 ORDER BY query_id;


-- name: CountQueries :one
SELECT COUNT(*) FROM mm_queries;
