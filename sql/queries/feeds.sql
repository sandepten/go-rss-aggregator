-- name: CreateFeed :one
INSERT INTO
  feeds (
    id,
    user_id,
    url,
    name,
    created_at,
    updated_at
  )
VALUES
  (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
  ) RETURNING *;

-- name: GetAllFeeds :many
SELECT
  *
FROM
  feeds
ORDER BY
  created_at DESC;

-- name: GetNextFeedsToFetch :many
SELECT
  *
FROM
  feeds
ORDER
  BY
    last_fetched_at ASC NULLS FIRST
LIMIT
  $1;

-- name: MarkFeedAsFetched :one
UPDATE
  feeds
SET
  last_fetched_at = NOW(),
  updated_at = NOW()
WHERE
  id = $1
RETURNING
  *;