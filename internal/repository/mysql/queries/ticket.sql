-- name: CreateTicket :execlastid
INSERT INTO mm_tickets (uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetTicketByID :one
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE ticket_id = ? LIMIT 1;

-- name: GetTicketByUUID :one
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE uuid = ? LIMIT 1;

-- name: GetTicketByCronID :one
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE cron_id = ? LIMIT 1;

-- name: UpdateTicket :exec
UPDATE mm_tickets SET uuid = ?, cluster_id = ?, `database` = ?, subject = ?, content = ?, status = ?, user_id = ?, reviewer_id = ?, cron_id = ?, version = ? WHERE ticket_id = ?;

-- name: DeleteTicketByID :exec
DELETE FROM mm_tickets WHERE ticket_id = ?;

-- name: DeleteTicketByUUID :exec
DELETE FROM mm_tickets WHERE uuid = ?;

-- name: ListTickets :many
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE ticket_id < ? ORDER BY ticket_id DESC LIMIT ?;

-- name: ListTicketsFiltered :many
-- status 传 0 表示不限状态；pattern 形如 %关键词%，空关键词用 %% 匹配全部
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets
WHERE ticket_id < ?
  AND (status = ? OR ? = 0)
  AND (subject LIKE ? OR `database` LIKE ?)
ORDER BY ticket_id DESC LIMIT ?;

-- name: ListTicketsByUser :many
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE user_id = ? ORDER BY ticket_id;

-- name: ListTicketsByReviewer :many
SELECT ticket_id, uuid, cluster_id, `database`, subject, content, status, user_id, reviewer_id, cron_id, version, update_at, create_at FROM mm_tickets WHERE reviewer_id = ? ORDER BY ticket_id;

-- name: CountTicketsByStatus :many
SELECT status, COUNT(*) AS count FROM mm_tickets WHERE create_at >= ? AND create_at < ? GROUP BY status ORDER BY status;

-- name: CountTicketsFiltered :one
-- 与 ListTicketsFiltered 同一口径；cursor 不参与计数，统计全部匹配行
SELECT COUNT(*) FROM mm_tickets
WHERE (status = ? OR ? = 0)
  AND (subject LIKE ? OR `database` LIKE ?);
