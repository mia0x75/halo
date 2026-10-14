-- name: CreateGlossary :exec
INSERT INTO mm_glossaries (`group`, `key`, value, uuid, description, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetGlossaryByUUID :one
SELECT `group`, `key`, value, uuid, description, version, update_at, create_at FROM mm_glossaries WHERE uuid = ? LIMIT 1;

-- name: GetGlossaryByKey :one
SELECT `group`, `key`, value, uuid, description, version, update_at, create_at FROM mm_glossaries WHERE `group` = ? AND `key` = ? LIMIT 1;

-- name: UpdateGlossary :exec
UPDATE mm_glossaries SET value = ?, description = ?, version = ?, update_at = ? WHERE uuid = ?;

-- name: DeleteGlossaryByUUID :exec
DELETE FROM mm_glossaries WHERE uuid = ?;

-- name: DeleteGlossaryByKey :exec
DELETE FROM mm_glossaries WHERE `group` = ? AND `key` = ?;

-- name: ListGlossaries :many
SELECT `group`, `key`, value, uuid, description, version, update_at, create_at FROM mm_glossaries WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;
