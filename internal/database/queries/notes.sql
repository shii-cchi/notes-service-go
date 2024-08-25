-- name: GetNotes :many
SELECT name, content FROM notes
WHERE user_id = $1;