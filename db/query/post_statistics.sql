-- name: IncrementLikes :exec
INSERT INTO post_statistics (post_id, likes, updated_at)
VALUES ($1, 1, CURRENT_TIMESTAMP)
    ON CONFLICT (post_id)
    DO UPDATE SET likes = post_statistics.likes + 1,
    updated_at = CURRENT_TIMESTAMP;

-- name: IncrementComments :exec
INSERT INTO post_statistics (post_id, comments, updated_at)
VALUES ($1, 1, CURRENT_TIMESTAMP)
    ON CONFLICT (post_id)
    DO UPDATE SET comments = post_statistics.comments + 1,
    updated_at = CURRENT_TIMESTAMP;
