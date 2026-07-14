-- name: IndexHandler :one
INSERT INTO link (
    SID, URL
) VALUES (
    ?, ?
)
RETURNING *;

-- name: SubmitHandler :one
UPDATE link
SET SID = ?,
    URL = ?
WHERE id = ?
RETURNING *;

-- name: RedirectHandler :one
SELECT URL
FROM link
WHERE SID = ? LIMIT 1;

-- name: CheckID :one
SELECT EXISTS(
    SELECT 1
    FROM link
    WHERE id = ?
    LIMIT 1
);
