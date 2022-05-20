package service

import (
	"context"
	"io"
	"library-go/internal/domain"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock.go

type ArticleService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Article, error)
	GetAll(sortingOptions *domain.SortFilterPagination) ([]*domain.Article, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, article *domain.CreateArticleDTO) (string, error)
	Update(ctx context.Context, article *domain.UpdateArticleDTO) error
	Load(ctx context.Context, path string) ([]byte, error)
	Save(ctx context.Context, path, fileName string, file io.Reader) error
}

type AudioService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Audio, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Audio, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, audio *domain.CreateAudioDTO) (string, error)
	Update(ctx context.Context, audio *domain.UpdateAudioDTO) error
	Load(ctx context.Context, path string) ([]byte, error)
	Save(ctx context.Context, path, fileName string, file io.Reader) error
}

type AuthorService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Author, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Author, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, author *domain.CreateAuthorDTO) (string, error)
	Update(ctx context.Context, author *domain.UpdateAuthorDTO) error
}

type BookService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Book, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Book, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, book *domain.CreateBookDTO) (string, error)
	Update(ctx context.Context, book *domain.UpdateBookDTO) error
	Load(ctx context.Context, path string) ([]byte, error)
	Save(ctx context.Context, path, fileName string, file io.Reader) error
}

type DirectionService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Direction, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Direction, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, direction *domain.CreateDirectionDTO) (string, error)
	Update(ctx context.Context, direction *domain.UpdateDirectionDTO) error
}

type ReviewService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Review, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Review, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, review *domain.CreateReviewDTO) (string, error)
	Update(ctx context.Context, review *domain.UpdateReviewDTO) error
}

type TagService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Tag, error)
	GetManyByUUIDs(ctx context.Context, UUIDs []string) ([]*domain.Tag, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Tag, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, tag *domain.CreateTagDTO) (string, error)
	Update(ctx context.Context, tag *domain.UpdateTagDTO) error
}

type VideoService interface {
	GetByUUID(ctx context.Context, UUID string) (*domain.Video, error)
	GetAll(ctx context.Context, limit, offset int) ([]*domain.Video, error)
	Delete(ctx context.Context, UUID string) error
	Create(ctx context.Context, video *domain.CreateVideoDTO) (string, error)
	Update(ctx context.Context, video *domain.UpdateVideoDTO) error
	Load(ctx context.Context, path string) ([]byte, error)
	Save(ctx context.Context, path, fileName string, file io.Reader) error
}
