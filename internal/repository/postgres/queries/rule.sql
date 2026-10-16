-- name: CreateRule :exec
INSERT INTO mm_rules (name, uuid, "group", description, level, vldr_group, operator, parameters, bitwise, func, message, element, version) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);

-- name: GetRuleByUUID :one
SELECT name, uuid, "group", description, level, vldr_group, operator, parameters, bitwise, func, message, element, version, update_at, create_at FROM mm_rules WHERE uuid = $1 LIMIT 1;

-- name: GetRuleByName :one
SELECT name, uuid, "group", description, level, vldr_group, operator, parameters, bitwise, func, message, element, version, update_at, create_at FROM mm_rules WHERE name = $1 LIMIT 1;

-- name: UpdateRule :exec
UPDATE mm_rules SET "group" = $1, description = $2, level = $3, vldr_group = $4, operator = $5, parameters = $6, bitwise = $7, func = $8, message = $9, element = $10, version = $11, update_at = NOW() WHERE uuid = $12;

-- name: DeleteRuleByUUID :exec
DELETE FROM mm_rules WHERE uuid = $1;

-- name: ListRules :many
SELECT name, uuid, "group", description, level, vldr_group, operator, parameters, bitwise, func, message, element, version, update_at, create_at FROM mm_rules WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;
