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

func (s *videoService) GetByUUID(UUID string) (*domain.Video, error) {
	return s.storage.GetOne(UUID)
}

func (s *videoService) GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Video, int, error) {
	return s.storage.GetAll(sortingOptions)
}

func (s *videoService) Delete(UUID, path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return s.storage.Delete(UUID)
}

func (s *videoService) Create(video *domain.CreateVideoDTO) (string, error) {
	return s.storage.Create(video)
}

func (s *videoService) Update(video *domain.UpdateVideoDTO) error {
	return s.storage.Update(video)
}

func (s *videoService) LoadFile(path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *videoService) SaveFile(path, fileName string, file io.Reader) error {
	return nil
	//return utils.SaveFile(path, fileName, file)
}

func (s *videoService) UpdateFile(dto *domain.UpdateVideoFileDTO) error {
	os.Remove(fmt.Sprintf("%s%s", dto.LocalPath, dto.OldFileName))

	err := s.storage.Update(&domain.UpdateVideoDTO{
		UUID:     dto.UUID,
		LocalURL: dto.LocalURL,
	})
	if err != nil {
		return err
	}

	return utils.SaveFile(dto.LocalPath, dto.NewFileName, dto.File)
}

func (s *videoService) Rate(UUID string, rating float32) error {
	return s.storage.Rate(UUID, rating)
}

func (s *videoService) DownloadCountUp(UUID string) error {
	return s.storage.DownloadCountUp(UUID)
}
