-- name: CreateBook :one
INSERT INTO books (title, author)
VALUES ($1, $2)
RETURNING id, title, author, created_at;

-- name: GetBook :one
SELECT *
FROM books
WHERE id = $1;

-- name: GetBooks :many
SELECT *
FROM books
ORDER BY created_at DESC;

-- name: UpdateBook :one
UPDATE books
SET title = $1, author = $2
WHERE id = $3
RETURNING id, title, author, created_at;

-- name: DeleteBook :exec
DELETE FROM books
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (email, password, username)
VALUES ($1, $2, $3)
RETURNING id, email, username, created_at;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;