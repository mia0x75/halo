-- name: CreateGlossary :exec
INSERT INTO mm_glossaries ("group", "key", value, uuid, description, version) VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetGlossaryByUUID :one
SELECT "group", "key", value, uuid, description, version, update_at, create_at FROM mm_glossaries WHERE uuid = $1 LIMIT 1;

-- name: GetGlossaryByKey :one
SELECT "group", "key", value, uuid, description, version, update_at, create_at FROM mm_glossaries WHERE "group" = $1 AND "key" = $2 LIMIT 1;

-- name: UpdateGlossary :exec
UPDATE mm_glossaries SET value = $1, description = $2, version = $3, update_at = NOW() WHERE uuid = $4;

-- name: DeleteGlossaryByUUID :exec
DELETE FROM mm_glossaries WHERE uuid = $1;

-- name: DeleteGlossaryByKey :exec
DELETE FROM mm_glossaries WHERE "group" = $1 AND "key" = $2;

-- name: ListGlossaries :many
SELECT "group", "key", value, uuid, description, version, update_at, create_at FROM mm_glossaries WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;

