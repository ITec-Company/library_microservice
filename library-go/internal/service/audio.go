package service

import (
	"context"
	"io"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
)

type audioService struct {
	logger  *logging.Logger
	storage store.AudioStorage
}

func NewAudioService(storage store.AudioStorage, logger *logging.Logger) AudioService {
	return &audioService{
		logger:  logger,
		storage: storage,
	}
}

func (s *audioService) GetByUUID(ctx context.Context, UUID string) (*domain.Audio, error) {
	return s.storage.GetOne(UUID)
}

func (s *audioService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Audio, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *audioService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *audioService) Create(ctx context.Context, audio *domain.CreateAudioDTO) (string, error) {
	return s.storage.Create(audio)
}

func (s *audioService) Update(ctx context.Context, audio *domain.UpdateAudioDTO) error {
	return s.storage.Update(audio)
}

func (s *audioService) Load(ctx context.Context, path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *audioService) Save(ctx context.Context, path, fileName string, file io.Reader) error {
	return utils.SaveFile(path, fileName, file)
}
