// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: gists.sql

package database

import (
	"context"
)

const countGists = `-- name: CountGists :one
SELECT COUNT(*) FROM gists
`

func (q *Queries) CountGists(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countGists)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createGist = `-- name: CreateGist :one
INSERT INTO gists (id, hash, language, result) VALUES ($1, $2, $3, $4)
RETURNING id, hash, result, language, created_at
`

type CreateGistParams struct {
	ID       string
	Hash     string
	Language string
	Result   string
}

func (q *Queries) CreateGist(ctx context.Context, arg CreateGistParams) (Gist, error) {
	row := q.db.QueryRow(ctx, createGist,
		arg.ID,
		arg.Hash,
		arg.Language,
		arg.Result,
	)
	var i Gist
	err := row.Scan(
		&i.ID,
		&i.Hash,
		&i.Result,
		&i.Language,
		&i.CreatedAt,
	)
	return i, err
}

const getGist = `-- name: GetGist :one
SELECT id, hash, result, language, created_at FROM gists
WHERE hash = $1 LIMIT 1
`

func (q *Queries) GetGist(ctx context.Context, hash string) (Gist, error) {
	row := q.db.QueryRow(ctx, getGist, hash)
	var i Gist
	err := row.Scan(
		&i.ID,
		&i.Hash,
		&i.Result,
		&i.Language,
		&i.CreatedAt,
	)
	return i, err
}

const listGists = `-- name: ListGists :many
SELECT id, hash, result, language, created_at FROM gists
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ListGistsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListGists(ctx context.Context, arg ListGistsParams) ([]Gist, error) {
	rows, err := q.db.Query(ctx, listGists, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Gist
	for rows.Next() {
		var i Gist
		if err := rows.Scan(
			&i.ID,
			&i.Hash,
			&i.Result,
			&i.Language,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
