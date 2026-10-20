-- name: CreateRuleHit :exec
INSERT INTO mm_violations (uuid, ticket_id, sequence, rule_name, level, message) VALUES ($1, $2, $3, $4, $5, $6);

-- name: ListRuleHitsByTicket :many
SELECT violation_id, uuid, ticket_id, sequence, rule_name, level, message, create_at FROM mm_violations WHERE ticket_id = $1 ORDER BY sequence, violation_id;

-- name: CountRuleHitsByLevel :many
SELECT level, COUNT(*) AS count FROM mm_violations WHERE create_at >= $1 AND create_at < $2 GROUP BY level ORDER BY level;

-- name: CountRuleHitsByRule :many
SELECT rule_name, COUNT(*) AS count FROM mm_violations WHERE create_at >= $1 AND create_at < $2 GROUP BY rule_name ORDER BY count DESC, rule_name ASC LIMIT $3;

-- name: ListRuleHits :many
SELECT violation_id, uuid, ticket_id, sequence, rule_name, level, message, create_at FROM mm_violations ORDER BY create_at DESC, violation_id DESC LIMIT $1;

-- name: DeleteRuleHitsBefore :execrows
DELETE FROM mm_violations WHERE create_at < $1;

-- name: DeleteRuleHitsByStatement :execrows
DELETE FROM mm_violations WHERE ticket_id = $1 AND sequence = $2;

-- name: DeleteRuleHitsByTicket :execrows
DELETE FROM mm_violations WHERE ticket_id = $1;
