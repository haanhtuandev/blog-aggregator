-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7

)
RETURNING *;

-- name: GetFeeds :many
SELECT f.name, f.url, u.name 
from feeds f join users u on f.user_id = u.id;

-- name: GetFeedByUrl :one
Select * from feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE feeds.id = $1;

-- name: GetNextFeedToFetch :one
SELECT * from feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;