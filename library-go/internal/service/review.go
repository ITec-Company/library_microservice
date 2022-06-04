package service

import (
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type reviewService struct {
	logger  *logging.Logger
	storage store.ReviewStorage
}

func NewReviewService(storage store.ReviewStorage, logger *logging.Logger) ReviewService {
	return &reviewService{
		logger:  logger,
		storage: storage}
}

func (s *reviewService) GetByUUID(UUID string) (*domain.Review, error) {
	return s.storage.GetOne(UUID)
}

func (s *reviewService) GetAll(limit, offset int) ([]*domain.Review, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *reviewService) Delete(UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *reviewService) Create(review *domain.CreateReviewDTO) (string, error) {
	return s.storage.Create(review)
}

func (s *reviewService) Update(review *domain.UpdateReviewDTO) error {
	return s.storage.Update(review)
}
