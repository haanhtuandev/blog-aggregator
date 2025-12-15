-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted_feed_follow.*, feeds.name as feed_name, users.name as user_name 
FROM inserted_feed_follow
INNER JOIN feeds on feeds.id = inserted_feed_follow.feed_id
INNER JOIN users on users.id = inserted_feed_follow.user_id;


-- name: GetFeedFollowsForUser :many
SELECT feeds.name as feed_name, users.name as user_name 
FROM feed_follows
INNER JOIN feeds on feeds.id = feed_follows.feed_id
INNER JOIN users on users.id = feed_follows.user_id
WHERE users.name = $1;

-- name: DeleteFollow :exec
DELETE FROM feed_follows
WHERE feed_follows.feed_id = $1 AND feed_follows.user_id = $2;

