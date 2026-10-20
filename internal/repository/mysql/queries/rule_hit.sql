-- name: CreateRuleHit :exec
INSERT INTO mm_violations (uuid, ticket_id, sequence, rule_name, level, message) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListRuleHitsByTicket :many
SELECT violation_id, uuid, ticket_id, sequence, rule_name, level, message, create_at FROM mm_violations WHERE ticket_id = ? ORDER BY sequence, violation_id;

-- name: CountRuleHitsByLevel :many
SELECT level, COUNT(*) AS count FROM mm_violations WHERE create_at >= ? AND create_at < ? GROUP BY level ORDER BY level;

-- name: CountRuleHitsByRule :many
SELECT rule_name, COUNT(*) AS count FROM mm_violations WHERE create_at >= ? AND create_at < ? GROUP BY rule_name ORDER BY count DESC, rule_name ASC LIMIT ?;

-- name: ListRuleHits :many
SELECT violation_id, uuid, ticket_id, sequence, rule_name, level, message, create_at FROM mm_violations ORDER BY create_at DESC, violation_id DESC LIMIT ?;
