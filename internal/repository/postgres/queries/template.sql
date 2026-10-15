-- name: CreateTemplate :exec
INSERT INTO mm_templates (uuid, subject, body, description, version) VALUES ($1, $2, $3, $4, $5);

-- name: GetTemplateByUUID :one
SELECT uuid, subject, body, description, version, update_at, create_at FROM mm_templates WHERE uuid = $1 LIMIT 1;

-- name: UpdateTemplate :exec
UPDATE mm_templates SET subject = $1, body = $2, description = $3, version = $4, update_at = NOW() WHERE uuid = $5;

-- name: DeleteTemplateByUUID :exec
DELETE FROM mm_templates WHERE uuid = $1;

-- name: ListTemplates :many
SELECT uuid, subject, body, description, version, update_at, create_at FROM mm_templates WHERE uuid < $1 ORDER BY uuid DESC LIMIT $2;

