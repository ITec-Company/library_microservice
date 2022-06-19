package service

import (
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type authorService struct {
	logger  *logging.Logger
	storage store.AuthorStorage
}

func NewAuthorService(storage *store.AuthorStorage, logger *logging.Logger) AuthorService {
	return &authorService{
		logger:  logger,
		storage: *storage}
}

func (s *authorService) GetByUUID(UUID string) (*domain.Author, error) {
	return s.storage.GetOne(UUID)
}

func (s *authorService) GetAll(limit, offset int) ([]*domain.Author, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *authorService) Delete(UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *authorService) Create(authorCreateDTO *domain.CreateAuthorDTO) (string, error) {
	return s.storage.Create(authorCreateDTO)

}

func (s *authorService) Update(author *domain.UpdateAuthorDTO) error {
	return s.storage.Update(author)
}
