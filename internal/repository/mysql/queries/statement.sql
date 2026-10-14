-- name: CreateStatement :exec
INSERT INTO mm_statements (ticket_id, sequence, uuid, content, `type`, status, report, plan, results, rows_affected, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetStatementByUUID :one
SELECT ticket_id, sequence, uuid, content, `type`, status, report, plan, results, rows_affected, version, update_at, create_at FROM mm_statements WHERE uuid = ? LIMIT 1;

-- name: GetStatementByTicketSequence :one
SELECT ticket_id, sequence, uuid, content, `type`, status, report, plan, results, rows_affected, version, update_at, create_at FROM mm_statements WHERE ticket_id = ? AND sequence = ? LIMIT 1;

-- name: UpdateStatement :exec
UPDATE mm_statements SET content = ?, `type` = ?, status = ?, report = ?, plan = ?, results = ?, rows_affected = ?, version = ?, update_at = ? WHERE uuid = ?;

-- name: DeleteStatementByUUID :exec
DELETE FROM mm_statements WHERE uuid = ?;

-- name: DeleteStatementByTicketSequence :exec
DELETE FROM mm_statements WHERE ticket_id = ? AND sequence = ?;

-- name: ListStatements :many
SELECT ticket_id, sequence, uuid, content, `type`, status, report, plan, results, rows_affected, version, update_at, create_at FROM mm_statements WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;
