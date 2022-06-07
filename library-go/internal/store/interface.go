package store

import (
	"library-go/internal/domain"
)

type ArticleStorage interface {
	GetOne(UUID string) (*domain.Article, error)
	GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Article, int, error)
	Create(article *domain.CreateArticleDTO) (string, error)
	Delete(UUID string) error
	Update(article *domain.UpdateArticleDTO) error
	Rate(UUID string, rating float32) error
	DownloadCountUp(UUID string) error
}

type AudioStorage interface {
	GetOne(UUID string) (*domain.Audio, error)
	GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Audio, int, error)
	Create(audio *domain.CreateAudioDTO) (string, error)
	Delete(UUID string) error
	Update(audio *domain.UpdateAudioDTO) error
	Rate(UUID string, rating float32) error
	DownloadCountUp(UUID string) error
}

type AuthorStorage interface {
	GetOne(UUID string) (*domain.Author, error)
	GetAll(limit, offset int) ([]*domain.Author, error)
	Create(author *domain.CreateAuthorDTO) (string, error)
	Delete(UUID string) error
	Update(author *domain.UpdateAuthorDTO) error
}

type BookStorage interface {
	GetOne(UUID string) (*domain.Book, error)
	GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Book, int, error)
	Create(book *domain.CreateBookDTO) (string, error)
	Delete(UUID string) error
	Update(book *domain.UpdateBookDTO) error
	Rate(UUID string, rating float32) error
	DownloadCountUp(UUID string) error
}

type DirectionStorage interface {
	GetOne(UUID string) (*domain.Direction, error)
	GetAll(limit, offset int) ([]*domain.Direction, error)
	Create(direction *domain.CreateDirectionDTO) (string, error)
	Delete(UUID string) error
	Update(direction *domain.UpdateDirectionDTO) error
}

type ReviewStorage interface {
	GetOne(UUID string) (*domain.Review, error)
	GetAll(limit, offset int) ([]*domain.Review, error)
	Create(review *domain.CreateReviewDTO) (string, error)
	Delete(UUID string) error
	Update(review *domain.UpdateReviewDTO) error
	Rate(UUID string, rating float32) error
}

type TagStorage interface {
	GetOne(UUID string) (*domain.Tag, error)
	GetMany(UUIDs []string) ([]*domain.Tag, error)
	GetAll(limit, offset int) ([]*domain.Tag, error)
	Create(tag *domain.CreateTagDTO) (string, error)
	Delete(UUID string) error
	Update(tag *domain.UpdateTagDTO) error
}

type VideoStorage interface {
	GetOne(UUID string) (*domain.Video, error)
	GetAll(sortOptions *domain.SortFilterPagination) ([]*domain.Video, int, error)
	Create(video *domain.CreateVideoDTO) (string, error)
	Delete(UUID string) error
	Update(video *domain.UpdateVideoDTO) error
	Rate(UUID string, rating float32) error
	DownloadCountUp(UUID string) error
}
