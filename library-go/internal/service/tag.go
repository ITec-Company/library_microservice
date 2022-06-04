package service

import (
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type tagService struct {
	logger  *logging.Logger
	storage store.TagStorage
}

func NewTagService(storage store.TagStorage, logger *logging.Logger) TagService {
	return &tagService{
		logger:  logger,
		storage: storage}
}

func (s *tagService) GetByUUID(UUID string) (*domain.Tag, error) {
	return s.storage.GetOne(UUID)
}

func (s *tagService) GetManyByUUIDs(UUIDs []string) ([]*domain.Tag, error) {
	return s.storage.GetMany(UUIDs)
}

func (s *tagService) GetAll(limit, offset int) ([]*domain.Tag, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *tagService) Delete(UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *tagService) Create(tag *domain.CreateTagDTO) (string, error) {
	return s.storage.Create(tag)
}

func (s *tagService) Update(tag *domain.UpdateTagDTO) error {
	return s.storage.Update(tag)
}
