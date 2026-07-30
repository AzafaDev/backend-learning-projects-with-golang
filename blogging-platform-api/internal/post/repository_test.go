package post

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_CreatePost(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	req := CreatePostRequest{
		Title:    "Test Title",
		Content:  "Test Content",
		Category: "Tech",
		Tags:     []string{"go", "testing"},
	}

	fixedID := uuid.New()
	now := time.Now()

	rows := pgxmock.NewRows([]string{"id", "title", "content", "category", "tags", "created_at", "updated_at"}).
		AddRow(fixedID, req.Title, req.Content, req.Category, req.Tags, now, now)

	mock.ExpectQuery("INSERT INTO posts").
		WithArgs(req.Title, req.Content, req.Category, req.Tags).
		WillReturnRows(rows)

	result, err := repo.CreatePost(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, fixedID, result.ID)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, req.Content, result.Content)
	assert.Equal(t, req.Category, result.Category)
	assert.Equal(t, req.Tags, result.Tags)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindPostByID_Found(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	id := uuid.New()
	now := time.Now()

	rows := pgxmock.NewRows([]string{"id", "title", "content", "category", "tags", "created_at", "updated_at"}).
		AddRow(id, "Title", "Content", "Tech", []string{"go"}, now, now)

	mock.ExpectQuery("SELECT (.+) FROM posts").
		WithArgs(id).
		WillReturnRows(rows)

	result, err := repo.FindPostByID(context.Background(), id)

	require.NoError(t, err)
	assert.Equal(t, id, result.ID)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_FindPostByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	id := uuid.New()

	mock.ExpectQuery("SELECT (.+) FROM posts").
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.FindPostByID(context.Background(), id)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrPostNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeletePostByID_NotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	id := uuid.New()

	mock.ExpectQuery("DELETE FROM posts").
		WithArgs(id).
		WillReturnError(pgx.ErrNoRows)

	result, err := repo.DeletePostByID(context.Background(), id)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrPostNotFound)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SearchPosts_WithTerm(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "title", "content", "category", "tags", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Tech Post", "About tech", "Technology", []string{"tech"}, now, now)

	mock.ExpectQuery("SELECT (.+) FROM posts WHERE").
		WithArgs("%tech%").
		WillReturnRows(rows)

	result, err := repo.SearchPosts(context.Background(), "tech")

	require.NoError(t, err)
	assert.Len(t, result, 1)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SearchPosts_EmptyTerm_ReturnsAll(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	now := time.Now()
	rows := pgxmock.NewRows([]string{"id", "title", "content", "category", "tags", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Post 1", "Content 1", "Tech", []string{}, now, now).
		AddRow(uuid.New(), "Post 2", "Content 2", "Life", []string{}, now, now)

	mock.ExpectQuery("SELECT (.+) FROM posts").
		WillReturnRows(rows)

	result, err := repo.SearchPosts(context.Background(), "")

	require.NoError(t, err)
	assert.Len(t, result, 2)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SearchPosts_NoResults_ReturnsEmptySlice(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &Repository{pool: mock}

	rows := pgxmock.NewRows([]string{"id", "title", "content", "category", "tags", "created_at", "updated_at"})

	mock.ExpectQuery("SELECT (.+) FROM posts").
		WillReturnRows(rows)

	result, err := repo.SearchPosts(context.Background(), "")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)

	require.NoError(t, mock.ExpectationsWereMet())
}
