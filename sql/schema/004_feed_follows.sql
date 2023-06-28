-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  unique(user_id, feed_id) -- this will prevent a user from following the same feed twice
);

-- +goose Down
DROP TABLE feed_follows;
