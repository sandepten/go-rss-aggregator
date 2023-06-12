-- name: CreateUser :one
INSERT INTO
  users (name, email, password, created_at, updated_at)
VALUES
  ($ 1, $ 2, $ 3, $ 4, $ 5) RETURNING id,
  name,
  email,
  created_at,
  updated_at;
