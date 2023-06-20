// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO
  users (id, email, password, created_at, updated_at)
VALUES
  ($1, $2, $3, $4, $5) RETURNING id, email, password, created_at, updated_at
`

type CreateUserParams struct {
	ID        uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
