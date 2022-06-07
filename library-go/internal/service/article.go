package service

import (
	"fmt"
	"image"
	"io"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
	"os"
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

func (s *articleService) GetByUUID(UUID string) (*domain.Article, error) {
	return s.storage.GetOne(UUID)
}

func (s *articleService) GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Article, int, error) {
	return s.storage.GetAll(sortingOptions)
}

func (s *articleService) Delete(UUID, path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return s.storage.Delete(UUID)
}

func (s *articleService) Create(article *domain.CreateArticleDTO) (string, error) {
	return s.storage.Create(article)
}

func (s *articleService) Update(article *domain.UpdateArticleDTO) error {
	return s.storage.Update(article)
}

func (s *articleService) LoadFile(path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *articleService) SaveFile(path, fileName string, file io.Reader) error {
	return utils.SaveFile(path, fileName, file)
}

func (s *articleService) UpdateFile(dto *domain.UpdateArticleFileDTO) error {
	if dto.OldFileName != dto.NewFileName {
		err := os.Remove(fmt.Sprintf("%s%s", dto.LocalPath, dto.OldFileName))
		if err != nil {
			return err
		}

		err = utils.SaveFile(dto.LocalPath, dto.NewFileName, dto.File)
		if err != nil {
			return err
		}

		return s.storage.Update(&domain.UpdateArticleDTO{
			UUID:     dto.UUID,
			LocalURL: dto.LocalURL,
		})

	} else {
		return utils.SaveFile(dto.LocalPath, dto.NewFileName, dto.File)
	}
}

func (s *articleService) LoadImage(path string, format utils.Format, extension utils.Extension) (*image.Image, error) {
	return utils.GetImageFromLocalStore(path, format, extension)
}

func (s *articleService) SaveImage(path string, image *image.Image) error {
	return utils.SaveImage(image, path)
}

func (s *articleService) UpdateImage(path string, image *image.Image) error {
	return utils.SaveImage(image, path)
}

func (s *articleService) Rate(UUID string, rating float32) error {
	return s.storage.Rate(UUID, rating)
}

func (s *articleService) DownloadCountUp(UUID string) error {
	return s.storage.DownloadCountUp(UUID)
}
