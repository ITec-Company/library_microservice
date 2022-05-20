package service

import (
	"context"
	"io"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
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

func (s *articleService) GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Article, int, error) {
	return s.storage.GetAll(sortingOptions)
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

func (s *articleService) Load(ctx context.Context, path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *articleService) Save(ctx context.Context, path, fileName string, file io.Reader) error {
	return utils.SaveFile(path, fileName, file)
}
