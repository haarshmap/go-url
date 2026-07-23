-- name: CreateUser :one
INSERT INTO users (
    username, hash_password, email
) VALUES (
    ?, ?, ?
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = ?
LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = ?
LIMIT 1;

-- name: CreateLink :one
INSERT INTO links (
    short_id, orig_url, expiry, user_id
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;

-- name: GetURLByShortCode :one
SELECT orig_url FROM links
WHERE short_id = ?
LIMIT 1;