-- name: CreateRule :exec
INSERT INTO mm_rules (name, uuid, `group`, description, level, vldr_group, operator, `values`, bitwise, func, message, element, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetRuleByUUID :one
SELECT name, uuid, `group`, description, level, vldr_group, operator, `values`, bitwise, func, message, element, version, update_at, create_at FROM mm_rules WHERE uuid = ? LIMIT 1;

-- name: GetRuleByName :one
SELECT name, uuid, `group`, description, level, vldr_group, operator, `values`, bitwise, func, message, element, version, update_at, create_at FROM mm_rules WHERE name = ? LIMIT 1;

-- name: UpdateRule :exec
UPDATE mm_rules SET `group` = ?, description = ?, level = ?, vldr_group = ?, operator = ?, `values` = ?, bitwise = ?, func = ?, message = ?, element = ?, version = ?, update_at = ? WHERE uuid = ?;

-- name: DeleteRuleByUUID :exec
DELETE FROM mm_rules WHERE uuid = ?;

-- name: ListRules :many
SELECT name, uuid, `group`, description, level, vldr_group, operator, `values`, bitwise, func, message, element, version, update_at, create_at FROM mm_rules WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;
