package service

import (
	"context"

	"vniizht/internal/model"
	"vniizht/internal/repository"
)

type SearchService struct {
	repo *repository.SearchRepo
}

func NewSearchService(repo *repository.SearchRepo) *SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(ctx context.Context, f model.SearchFilter) ([]*model.SearchResult, error) {
	return s.repo.Search(ctx, f)
}
