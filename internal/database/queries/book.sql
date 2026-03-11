-- name: CreateBook :exec
INSERT INTO books (id, title, description, author_id)
VALUES (?, ?, ?, ?);

-- name: GetAllBooks :many
SELECT
    books.id,
    books.title,
    books.description,
    authors.id   AS author_id,
    authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
ORDER BY books.id;

-- name: GetBookByID :one
SELECT
    books.id,
    books.title,
    books.description,
    authors.id   AS author_id,
    authors.name AS author_name
FROM books
JOIN authors ON authors.id = books.author_id
WHERE books.id = ?;

-- name: UpdateBook :exec
UPDATE books
SET title = ?, description = ?
WHERE id = ?;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = ?;

-- name: CountBooks :one
SELECT count(*) FROM books;
