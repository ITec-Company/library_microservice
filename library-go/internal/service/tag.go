package service

import (
	"context"
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

func (s *tagService) GetByUUID(ctx context.Context, UUID string) (*domain.Tag, error) {
	return s.storage.GetOne(UUID)
}

func (s *tagService) GetManyByUUIDs(ctx context.Context, UUIDs []string) ([]*domain.Tag, error) {
	return s.storage.GetMany(UUIDs)
}

func (s *tagService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Tag, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *tagService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *tagService) Create(ctx context.Context, tag *domain.CreateTagDTO) (string, error) {
	return s.storage.Create(tag)
}

func (s *tagService) Update(ctx context.Context, tag *domain.UpdateTagDTO) error {
	return s.storage.Update(tag)
}
