-- name: CreateNote :one
INSERT INTO notes (name, content, user_id)
VALUES ($1, $2, $3)
RETURNING id, name, content;

-- name: GetNotes :many
SELECT id, name, content FROM notes
WHERE user_id = $1;