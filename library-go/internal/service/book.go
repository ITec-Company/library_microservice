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

type bookService struct {
	logger  *logging.Logger
	storage store.BookStorage
}

func NewBookService(storage store.BookStorage, logger *logging.Logger) BookService {
	return &bookService{
		logger:  logger,
		storage: storage,
	}
}

func (s *bookService) GetByUUID(UUID string) (*domain.Book, error) {
	return s.storage.GetOne(UUID)
}

func (s *bookService) GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Book, int, error) {
	return s.storage.GetAll(sortingOptions)
}

func (s *bookService) Delete(UUID, path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return s.storage.Delete(UUID)
}

func (s *bookService) Create(book *domain.CreateBookDTO) (string, error) {
	return s.storage.Create(book)
}

func (s *bookService) Update(book *domain.UpdateBookDTO) error {
	return s.storage.Update(book)
}

func (s *bookService) LoadFile(path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *bookService) SaveFile(path, fileName string, file io.Reader) error {
	return utils.SaveFile(path, fileName, file)
}

func (s *bookService) UpdateFile(dto *domain.UpdateBookFileDTO) error {
	if dto.OldFileName != dto.NewFileName {
		err := os.Remove(fmt.Sprintf("%s%s", dto.LocalPath, dto.OldFileName))
		if err != nil {
			return err
		}

		err = utils.SaveFile(dto.LocalPath, dto.NewFileName, dto.File)
		if err != nil {
			return err
		}

		return s.storage.Update(&domain.UpdateBookDTO{
			UUID:     dto.UUID,
			LocalURL: dto.LocalURL,
		})

	} else {
		return utils.SaveFile(dto.LocalPath, dto.NewFileName, dto.File)
	}
}

func (s *bookService) LoadImage(path string, format utils.Format, extension utils.Extension) (*image.Image, error) {
	return utils.GetImageFromLocalStore(path, format, extension)
}

func (s *bookService) SaveImage(path string, image *image.Image) error {
	return utils.SaveImage(image, path)
}

func (s *bookService) UpdateImage(path string, image *image.Image) error {
	return utils.SaveImage(image, path)
}
