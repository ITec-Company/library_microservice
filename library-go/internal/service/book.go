package service

import (
	"context"
	"io"
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
	"library-go/pkg/utils"
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

func (s *bookService) GetByUUID(ctx context.Context, UUID string) (*domain.Book, error) {
	return s.storage.GetOne(UUID)
}

func (s *bookService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Book, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *bookService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *bookService) Create(ctx context.Context, book *domain.CreateBookDTO) (string, error) {
	return s.storage.Create(book)
}

func (s *bookService) Update(ctx context.Context, book *domain.UpdateBookDTO) error {
	return s.storage.Update(book)
}

func (s *bookService) Load(ctx context.Context, path string) ([]byte, error) {
	return utils.LoadFile(path)
}

func (s *bookService) Save(ctx context.Context, path, fileName string, file io.Reader) error {
	return utils.SaveFile(path, fileName, file)
}
