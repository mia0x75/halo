-- name: CreateTicket :one
INSERT INTO mm_tickets (uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING ticket_id;

-- name: GetTicketByID :one
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE ticket_id = $1 LIMIT 1;

-- name: GetTicketByUUID :one
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE uuid = $1 LIMIT 1;

-- name: GetTicketByCronID :one
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE cron_id = $1 LIMIT 1;

-- name: UpdateTicket :exec
UPDATE mm_tickets SET uuid = $1, cluster_id = $2, "database" = $3, subject = $4, content = $5, status = $6, user_id = $7, reviewer_id = $8, cron_id = $9, version = $10, update_at = NOW() WHERE ticket_id = $11;

-- name: DeleteTicketByID :exec
DELETE FROM mm_tickets WHERE ticket_id = $1;

-- name: DeleteTicketByUUID :exec
DELETE FROM mm_tickets WHERE uuid = $1;

-- name: ListTickets :many
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE ticket_id < $1 ORDER BY ticket_id DESC LIMIT $2;

-- name: ListTicketsByUser :many
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE user_id = $1 ORDER BY ticket_id;


-- name: CountTicketsByStatus :many
SELECT status, COUNT(*) AS count FROM mm_tickets WHERE create_at >= $1 AND create_at < $2 GROUP BY status ORDER BY status;
