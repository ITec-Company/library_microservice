package service

import (
	"image"
	"io"
	"library-go/internal/domain"
	"library-go/pkg/utils"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type ArticleService interface {
	GetByUUID(UUID string) (*domain.Article, error)
	GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Article, int, error)
	Delete(UUID, path string) error
	Create(article *domain.CreateArticleDTO) (string, error)
	Update(article *domain.UpdateArticleDTO) error
	LoadFile(path string) ([]byte, error)
	SaveFile(path, fileName string, file io.Reader) error
	UpdateFile(dto *domain.UpdateArticleFileDTO) error
	LoadImage(path string, format utils.Format, extension utils.Extension) (*image.Image, error)
	SaveImage(path string, image *image.Image) error
	UpdateImage(path string, image *image.Image) error
}

type AudioService interface {
	GetByUUID(UUID string) (*domain.Audio, error)
	GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Audio, int, error)
	Delete(UUID, path string) error
	Create(audio *domain.CreateAudioDTO) (string, error)
	Update(audio *domain.UpdateAudioDTO) error
	LoadFile(path string) ([]byte, error)
	SaveFile(path, fileName string, file io.Reader) error
	UpdateFile(dto *domain.UpdateAudioFileDTO) error
}

type AuthorService interface {
	GetByUUID(UUID string) (*domain.Author, error)
	GetAll(limit, offset int) ([]*domain.Author, error)
	Delete(UUID string) error
	Create(author *domain.CreateAuthorDTO) (string, error)
	Update(author *domain.UpdateAuthorDTO) error
}

type BookService interface {
	GetByUUID(UUID string) (*domain.Book, error)
	GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Book, int, error)
	Delete(UUID, path string) error
	Create(book *domain.CreateBookDTO) (string, error)
	Update(book *domain.UpdateBookDTO) error
	LoadFile(path string) ([]byte, error)
	SaveFile(path, fileName string, file io.Reader) error
	UpdateFile(dto *domain.UpdateBookFileDTO) error
	LoadImage(path string, format utils.Format, extension utils.Extension) (*image.Image, error)
	SaveImage(path string, image *image.Image) error
	UpdateImage(path string, image *image.Image) error
}

type DirectionService interface {
	GetByUUID(UUID string) (*domain.Direction, error)
	GetAll(limit, offset int) ([]*domain.Direction, error)
	Delete(UUID string) error
	Create(direction *domain.CreateDirectionDTO) (string, error)
	Update(direction *domain.UpdateDirectionDTO) error
}

type ReviewService interface {
	GetByUUID(UUID string) (*domain.Review, error)
	GetAll(limit, offset int) ([]*domain.Review, error)
	Delete(UUID string) error
	Create(review *domain.CreateReviewDTO) (string, error)
	Update(review *domain.UpdateReviewDTO) error
}

type TagService interface {
	GetByUUID(UUID string) (*domain.Tag, error)
	GetManyByUUIDs(UUIDs []string) ([]*domain.Tag, error)
	GetAll(limit, offset int) ([]*domain.Tag, error)
	Delete(UUID string) error
	Create(tag *domain.CreateTagDTO) (string, error)
	Update(tag *domain.UpdateTagDTO) error
}

type VideoService interface {
	GetByUUID(UUID string) (*domain.Video, error)
	GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Video, int, error)
	Delete(UUID, path string) error
	Create(video *domain.CreateVideoDTO) (string, error)
	Update(video *domain.UpdateVideoDTO) error
	LoadFile(path string) ([]byte, error)
	SaveFile(path, fileName string, file io.Reader) error
	UpdateFile(dto *domain.UpdateVideoFileDTO) error
}
