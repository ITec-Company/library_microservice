package service

import (
	"fmt"
	"io"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
	"os"
)

type audioService struct {
	logger  *logging.Logger
	storage store.AudioStorage
}

func NewAudioService(storage *store.AudioStorage, logger *logging.Logger) AudioService {
	return &audioService{
		logger:  logger,
		storage: *storage,
	}
}

func (s *audioService) GetByUUID(UUID string) (*domain.Audio, error) {
	return s.storage.GetOne(UUID)
}

func (s *audioService) GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Audio, int, error) {
	return s.storage.GetAll(sortingOptions)
}

func (s *audioService) Delete(UUID, path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return s.storage.Delete(UUID)
}

func (s *audioService) Create(audio *domain.CreateAudioDTO) (string, error) {
	return s.storage.Create(audio)
}

func (s *audioService) Update(audio *domain.UpdateAudioDTO) error {
	return s.storage.Update(audio)
}

func (s *audioService) LoadFile(path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *audioService) SaveFile(path, fileName string, file io.Reader) error {
	return utils.SaveFile(path, fileName, file)
}

func (s *audioService) UpdateFile(dto *domain.UpdateAudioFileDTO) error {
	os.Remove(fmt.Sprintf("%s%s", dto.LocalPath, dto.OldFileName))

	err := s.storage.Update(&domain.UpdateAudioDTO{
		UUID:     dto.UUID,
		LocalURL: dto.LocalURL,
	})
	if err != nil {
		return err
	}

	return utils.SaveFile(dto.LocalPath, dto.NewFileName, dto.File)
}

func (s *audioService) Rate(UUID string, rating float32) error {
	return s.storage.Rate(UUID, rating)
}

func (s *audioService) DownloadCountUp(UUID string) error {
	return s.storage.DownloadCountUp(UUID)
}
