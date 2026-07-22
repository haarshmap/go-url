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

-- name: CreateLink :one
INSERT INTO links (
    short_id, orig_url, expiry, user_id
) VALUES (
    ?, ?, ?, ?
)
RETURNING *;