package post

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) CreatePost(ctx context.Context, req CreatePostRequest) (*Post, error) {
	query := `
	INSERT INTO posts (title, content, category, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, title, content, category, tags, created_at, updated_at
	`

	var p Post
	err := r.pool.QueryRow(ctx, query,
		req.Title,
		req.Content,
		req.Category,
		req.Tags,
	).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		&p.Tags,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Repository) FindPostByID(ctx context.Context, id uuid.UUID) (*Post, error) {
	query := `
	SELECT id, title, content, category, tags, created_at, updated_at
	FROM posts
	WHERE id=$1
	`
	var p Post
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		&p.Tags,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (r *Repository) UpdatePostByID(ctx context.Context, req UpdatePostRequest, id uuid.UUID) (*Post, error) {
	query := `
	UPDATE posts
	SET title=$1,
	content=$2,
	category=$3,
	tags=$4
	WHERE id=$5
	RETURNING id, title, content, category, tags, created_at, updated_at
	`
	var p Post
	if err := r.pool.QueryRow(ctx, query,
		req.Title,
		req.Content,
		req.Category,
		req.Tags,
		id,
	).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		&p.Tags,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (r *Repository) DeletePostByID(ctx context.Context, id uuid.UUID) (*Post, error) {
	query := `
	DELETE FROM posts
	WHERE id=$1
	RETURNING id, title, content, category, tags, created_at, updated_at
	`
	var p Post
	if err := r.pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Category,
		&p.Tags,
		&p.CreatedAt,
		&p.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrPostNotFound
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func (r *Repository) SearchPosts(ctx context.Context, term string) ([]Post, error) {
	query := `
	SELECT id, title, content, category, tags, created_at, updated_at
	FROM posts
	`
	var args []any

	if term != "" {
		query += `
		WHERE title ILIKE $1
		   OR content ILIKE $1
		   OR category ILIKE $1
		`
		args = append(args, "%"+term+"%")
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ps []Post

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.Category,
			&p.Tags,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if ps == nil {
		return []Post{}, nil
	}

	return ps, nil
}
