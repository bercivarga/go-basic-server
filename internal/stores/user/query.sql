-- user/query.sql
-- ------------------------------------------------------------
-- Users basic queries for sqlc (SQLite engine)
-- Schema: id INTEGER PK, email TEXT UNIQUE, password_hash TEXT, created_at DATETIME
-- ------------------------------------------------------------

-- Create a new user and return the generated row --------------------------------
-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (?, ?)
RETURNING id, email, created_at;

-- Fetch a user by primary key ----------------------------------------------------
-- name: GetUserByID :one
SELECT id, email, password_hash, role, created_at
FROM   users
WHERE  id = ?;

-- Fetch a user by unique email ---------------------------------------------------
-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at
FROM   users
WHERE  email = ?;

-- Get user role -----------------------------------------------------------------
-- name: GetRole :one
SELECT role FROM users
WHERE id = ?;

-- List active users (simple pagination) -----------------------------------------
-- name: ListUsers :many
SELECT id, email, role, created_at
FROM   users
ORDER  BY id
LIMIT  ?  OFFSET ?;

-- Update only the password hash --------------------------------------------------
-- name: UpdatePasswordHash :exec
UPDATE users
SET    password_hash = ?
WHERE  id = ?;

-- Delete a user -----------------------------------------------------------------
-- name: DeleteUser :exec
DELETE FROM users
WHERE  id = ?;
