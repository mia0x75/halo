-- name: CreateTemplate :exec
INSERT INTO mm_templates (uuid, subject, body, description, version, update_at, create_at) VALUES (?, ?, ?, ?, ?, ?, ?);

-- name: GetTemplateByUUID :one
SELECT uuid, subject, body, description, version, update_at, create_at FROM mm_templates WHERE uuid = ? LIMIT 1;

-- name: UpdateTemplate :exec
UPDATE mm_templates SET subject = ?, body = ?, description = ?, version = ?, update_at = ? WHERE uuid = ?;

-- name: DeleteTemplateByUUID :exec
DELETE FROM mm_templates WHERE uuid = ?;

-- name: ListTemplates :many
SELECT uuid, subject, body, description, version, update_at, create_at FROM mm_templates WHERE uuid < ? ORDER BY uuid DESC LIMIT ?;
