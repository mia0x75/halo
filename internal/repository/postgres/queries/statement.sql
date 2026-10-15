-- name: CreateStatement :exec
INSERT INTO mm_statements (ticket_id, sequence, uuid, content, "type", status, report, plan, results, rows_affected, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);

-- name: GetStatementByUUID :one
SELECT ticket_id, sequence, uuid, content, "type", status, report, plan, results, rows_affected, version, update_at, create_at FROM mm_statements WHERE uuid = $1 LIMIT 1;

-- name: GetStatementByTicketSequence :one
SELECT ticket_id, sequence, uuid, content, "type", status, report, plan, results, rows_affected, version, update_at, create_at FROM mm_statements WHERE ticket_id = $1 AND sequence = $2 LIMIT 1;

-- name: UpdateStatement :exec
UPDATE mm_statements SET content = $1, "type" = $2, status = $3, report = $4, plan = $5, results = $6, rows_affected = $7, version = $8, update_at = NOW() WHERE uuid = $9;

-- name: DeleteStatementByUUID :exec
DELETE FROM mm_statements WHERE uuid = $1;

-- name: DeleteStatementByTicketSequence :exec
DELETE FROM mm_statements WHERE ticket_id = $1 AND sequence = $2;

-- name: DeleteStatementsByTicket :exec
DELETE FROM mm_statements WHERE ticket_id = $1;

-- name: ListStatements :many
SELECT ticket_id, sequence, uuid, content, "type", status, report, plan, results, rows_affected, version, update_at, create_at FROM mm_statements WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;

