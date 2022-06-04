package service

import (
	"library-go/internal/domain"
	"library-go/internal/store"
	"library-go/pkg/logging"
)

type directionService struct {
	logger  *logging.Logger
	storage store.DirectionStorage
}

func NewDirectionService(storage store.DirectionStorage, logger *logging.Logger) DirectionService {
	return &directionService{
		logger:  logger,
		storage: storage}
}

func (s *directionService) GetByUUID(UUID string) (*domain.Direction, error) {
	return s.storage.GetOne(UUID)
}

func (s *directionService) GetAll(limit, offset int) ([]*domain.Direction, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *directionService) Delete(UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *directionService) Create(direction *domain.CreateDirectionDTO) (string, error) {
	return s.storage.Create(direction)
}

func (s *directionService) Update(direction *domain.UpdateDirectionDTO) error {
	return s.storage.Update(direction)
}
