package service

import (
	"context"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type articleService struct {
	logger  *logging.Logger
	storage store.ArticleStorage
}

func NewArticleService(storage store.ArticleStorage, logger *logging.Logger) ArticleService {
	return &articleService{
		logger:  logger,
		storage: storage,
	}
}

func (s *articleService) GetByUUID(ctx context.Context, UUID string) (*domain.Article, error) {
	return s.storage.GetOne(UUID)
}

func (s *articleService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Article, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *articleService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *articleService) Create(ctx context.Context, article *domain.CreateArticleDTO) (string, error) {
	return s.storage.Create(article)
}

func (s *articleService) Update(ctx context.Context, article *domain.UpdateArticleDTO) error {
	return s.storage.Update(article)
}
