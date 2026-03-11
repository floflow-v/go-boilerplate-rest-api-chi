-- name: CreateAuthor :exec
INSERT INTO authors (id, name)
VALUES (?, ?);

-- name: GetAuthorByID :one
SELECT *
FROM authors
WHERE id = ?
