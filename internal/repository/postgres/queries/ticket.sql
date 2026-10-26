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

-- name: ListTicketsFiltered :many
-- status 传 0 表示不限状态；pattern 形如 %关键词%，空关键词用 %% 匹配全部
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets
WHERE ticket_id < $1
  AND (status = $2 OR $3 = 0)
  AND (subject LIKE $4 OR "database" LIKE $4)
ORDER BY ticket_id DESC LIMIT $5;

-- name: ListTicketsByUser :many
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE user_id = $1 ORDER BY ticket_id;


-- name: ListTicketsByReviewer :many
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE reviewer_id = $1 ORDER BY ticket_id;

-- name: ListExecutionTickets :many
-- status 由调用方固定为 LGTM；reviewer_id 传 0 表示不限审核人；scheduled 传 1 只看已预约
SELECT ticket_id, uuid, cluster_id, "database", subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets
WHERE ticket_id < $1
  AND status = $2
  AND (reviewer_id = $3 OR $4::bigint = 0)
  AND ($5::int = 0 OR cron_id IS NOT NULL)
ORDER BY ticket_id DESC LIMIT $6;

-- name: CountExecutionTickets :one
-- 与 ListExecutionTickets 同一口径；cursor 不参与计数
SELECT COUNT(*) FROM mm_tickets
WHERE status = $1
  AND (reviewer_id = $2 OR $3::bigint = 0)
  AND ($4::int = 0 OR cron_id IS NOT NULL);

-- name: CountTicketsByStatus :many
SELECT status, COUNT(*) AS count FROM mm_tickets WHERE create_at >= $1 AND create_at < $2 GROUP BY status ORDER BY status;

-- name: CountTicketsFiltered :one
-- 与 ListTicketsFiltered 同一口径；cursor 不参与计数，统计全部匹配行
SELECT COUNT(*) FROM mm_tickets
WHERE (status = $1 OR $2 = 0)
  AND (subject LIKE $3 OR "database" LIKE $3);
