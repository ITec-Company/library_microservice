package service

import (
	"context"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type videoService struct {
	logger  *logging.Logger
	storage store.VideoStorage
}

func NewService(storage store.VideoStorage, logger *logging.Logger) VideoService {
	return &videoService{
		logger:  logger,
		storage: storage,
	}
}

func (s *videoService) GetByUUID(ctx context.Context, UUID string) (*domain.Video, error) {
	return s.storage.GetOne(UUID)
}

func (s *videoService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Video, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *videoService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *videoService) Create(ctx context.Context, video *domain.CreateVideoDTO) (string, error) {
	return s.storage.Create(video)
}

func (s *videoService) Update(ctx context.Context, video *domain.UpdateVideoDTO) error {
	return s.storage.Update(video)
}
