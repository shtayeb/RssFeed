// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users.sql

package database

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, username,name,email,password)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, username, name, email, password
`

type CreateUserParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Name      string
	Email     string
	Password  string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Username,
		arg.Name,
		arg.Email,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getUserByAPIKey = `-- name: GetUserByAPIKey :one
SELECT id, created_at, updated_at, username, name, email, password FROM users WHERE id = $1
`

func (q *Queries) GetUserByAPIKey(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByAPIKey, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}

const getUserByEmailOrUsername = `-- name: GetUserByEmailOrUsername :one
SELECT id, created_at, updated_at, username, name, email, password FROM users WHERE email=$1 or username=$1
`

func (q *Queries) GetUserByEmailOrUsername(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailOrUsername, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Username,
		&i.Name,
		&i.Email,
		&i.Password,
	)
	return i, err
}
