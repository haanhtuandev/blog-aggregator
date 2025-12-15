-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
);

-- name: GetPostForUser :many
SELECT * from posts p
JOIN feeds f on p.feed_id = f.id
JOIN users u on u.id = f.user_id
WHERE user_id = $1
ORDER BY published_at DESC
LIMIT $2;