package post

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service struct {
	repo     *Repository
	validate *validator.Validate
}

func NewService(r *Repository) *Service {
	return &Service{
		repo:     r,
		validate: validator.New(),
	}
}

func (s *Service) CreatePost(ctx context.Context, req CreatePostRequest) (*Post, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, ErrBadRequest
	}

	return s.repo.CreatePost(ctx, req)
}

func (s *Service) GetPostByID(ctx context.Context, id uuid.UUID) (*Post, error) {
	return s.repo.FindPostByID(ctx, id)
}

func (s *Service) UpdatePost(ctx context.Context, req UpdatePostRequest, id uuid.UUID) (*Post, error) {
	if err := s.validate.Struct(req); err != nil {
		return nil, ErrBadRequest
	}
	return s.repo.UpdatePostByID(ctx, req, id)
}

func (s *Service) DeletePostByID(ctx context.Context, id uuid.UUID) (*Post, error) {
	return s.repo.DeletePostByID(ctx, id)
}

func (s *Service) SearchPosts(ctx context.Context, term string) ([]Post, error) {
	return s.repo.SearchPosts(ctx, term)
}
